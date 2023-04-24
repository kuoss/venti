package alerter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/store"
	"github.com/kuoss/venti/pkg/store/remote"
	commonModel "github.com/prometheus/common/model"
	promRule "github.com/prometheus/prometheus/rules"
)

type alerter struct {
	alertFiles         []model.AlertFile
	evaluationInterval time.Duration
	remoteStore        *remote.RemoteStore
	repeat             bool
	alertmanagerURL    string
}

func New(stores *store.Stores) *alerter {
	return &alerter{
		alertFiles:         getAlertFiles(stores),
		remoteStore:        stores.RemoteStore,
		alertmanagerURL:    "http://alertmanager:9093", // TODO: configurable
		evaluationInterval: 20 * time.Second,           // TODO: configurable
	}
}

func getAlertFiles(stores *store.Stores) []model.AlertFile {
	alertFiles := []model.AlertFile{}
	for _, ruleFile := range stores.AlertRuleStore.AlertRuleFiles() {
		datasources := stores.DatasourceStore.GetDatasourcesWithSelector(ruleFile.DatasourceSelector)
		if len(datasources) < 1 {
			logger.Warnf("no datasources from GetDatasourcesWithSelector")
			continue
		}
		for _, datasource := range datasources {
			alertGroups := []model.AlertGroup{}
			for _, ruleGroup := range ruleFile.RuleGroups {
				alerts := []model.Alert{}
				for _, rule := range ruleGroup.Rules {
					labels := ruleFile.CommonLabels
					for k, v := range rule.Labels {
						labels[k] = v
					}
					alerts = append(alerts, model.Alert{
						State:       promRule.StateInactive,
						Name:        rule.Alert,
						Expr:        rule.Expr,
						For:         rule.For,
						Labels:      labels,
						Annotations: rule.Annotations,
					})
				}
				alertGroups = append(alertGroups, model.AlertGroup{Alerts: alerts})
			}
			alertFiles = append(alertFiles, model.AlertFile{
				AlertGroups: alertGroups,
				Datasource:  datasource,
			})
		}
	}
	return alertFiles
}

func (a *alerter) SetAlertmanagerURL(url string) {
	a.alertmanagerURL = url
}

func (a *alerter) Start() {
	logger.Infof("starting alerter...")
	a.repeat = true
	go a.loop()
}

func (a *alerter) Stop() {
	logger.Infof("stopping alerter...")
	a.repeat = false
}

func (a *alerter) loop() {
	for {
		// TODO: test cover: go test -race
		if !a.repeat {
			logger.Infof("alerter stopped")
			return
		}
		a.Once()
		logger.Infof("sleep: %s", a.evaluationInterval)
		time.Sleep(a.evaluationInterval)
	}
}

func (a *alerter) Once() {
	logger.Infof("processAlertFiles")
	err := a.processAlertFiles()
	if err != nil {
		logger.Errorf("error on processAlertFiles: %s", err)
	}
}

func (a *alerter) processAlertFiles() error {
	totalFires := []model.Fire{}
	for i := range a.alertFiles {
		for j := range a.alertFiles[i].AlertGroups {
			for k := range a.alertFiles[i].AlertGroups[j].Alerts {
				fires, err := a.processAlert(&a.alertFiles[i].AlertGroups[j].Alerts[k], &a.alertFiles[i].Datasource)
				if err != nil {
					// TODO: test cover
					logger.Warnf("error on processAlert: %s", err)
					continue
				}
				totalFires = append(totalFires, fires...)
				time.Sleep(time.Duration(500) * time.Millisecond)
			}
		}
	}
	err := a.sendFires(totalFires)
	if err != nil {
		return fmt.Errorf("error on sendFires: %w", err)
	}
	return nil
}

func (a *alerter) processAlert(alert *model.Alert, datasource *model.Datasource) ([]model.Fire, error) {
	var zero []model.Fire
	queryData, err := a.queryAlert(alert, datasource)
	if err != nil {
		return zero, fmt.Errorf("error on queryAlert: %s", err)
	}
	fires := evaluateAlert(alert, queryData)
	return fires, nil
}

func (a *alerter) queryAlert(alert *model.Alert, datasource *model.Datasource) (model.QueryData, error) {
	var zero model.QueryData
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	// _ => queryResult.Status
	code, body, err := a.remoteStore.GET(ctx, datasource, "query", "query="+url.QueryEscape(alert.Expr))
	if err != nil {
		return zero, fmt.Errorf("error on GET: %w", err)
	}

	//  body: {"status":"success","data":{"resultType":"vector","result":[]}}
	var queryResult model.QueryResult
	err = json.Unmarshal([]byte(body), &queryResult)
	if err != nil {
		return zero, fmt.Errorf("error on Unmarshal: %w", err)
	}

	// wrap up
	if queryResult.Status != "success" {
		err = fmt.Errorf("not success status (status=%s, code=%d)", queryResult.Status, code)
	} else {
		if code != http.StatusOK {
			// maybe not reachable
			err = fmt.Errorf("not ok (code=%d)", code)
		}
	}
	return queryResult.Data, err
}

func evaluateAlert(alert *model.Alert, queryData model.QueryData) []model.Fire {
	var zero []model.Fire

	// inactive
	if len(queryData.Result) < 1 {
		alert.State = promRule.StateInactive
		alert.ActiveAt = 0
		return zero
	}

	if alert.ActiveAt == 0 {
		// now active
		alert.ActiveAt = commonModel.Now()
	}
	// pending
	if alert.ActiveAt.Add(time.Duration(alert.For)).After(commonModel.Now()) {
		alert.State = promRule.StatePending
		return zero
	}
	// firing
	alert.State = promRule.StateFiring
	return getFires(alert, queryData)
}

func getFires(alert *model.Alert, data model.QueryData) []model.Fire {
	if len(alert.Labels) < 1 {
		alert.Labels = map[string]string{}
	}
	if len(alert.Annotations) < 1 {
		alert.Annotations = map[string]string{}
	}
	alert.Labels["alertname"] = alert.Name
	if alert.Labels["alertname"] == "" {
		alert.Labels["alertname"] = "placeholder name"
	}
	alert.Labels["firer"] = "venti"
	if _, exists := alert.Annotations["summary"]; !exists {
		if len(alert.Annotations) == 0 {
			alert.Annotations = map[string]string{}
		}
		alert.Annotations["summary"] = "placeholder summary"
	}
	if data.ResultType != commonModel.ValVector {
		logger.Warnf("resultType is not vector")
		return []model.Fire{{
			State:       "firing",
			Labels:      alert.Labels,
			Annotations: alert.Annotations,
		}}
	}
	fires := []model.Fire{}
	for _, sample := range data.Result {
		fire := model.Fire{
			State:       "firing",
			Labels:      map[string]string{},
			Annotations: map[string]string{},
		}
		for k, v := range alert.Labels {
			fire.Labels[k] = v
		}
		for k, v := range alert.Annotations {
			fire.Annotations[k] = v
		}
		summary, err := renderSummary(alert.Annotations["summary"], &sample)
		if err != nil {
			logger.Warnf("error on renderSummary: %s", err)
		}
		fire.Annotations["summary"] = summary
		fires = append(fires, fire)
	}
	return fires
}

func renderSummary(input string, sample *commonModel.Sample) (string, error) {
	// pre-render
	text := input
	text = strings.ReplaceAll(text, "$value", sample.Value.String())
	text = strings.ReplaceAll(text, "$labels.", ".")
	text = strings.ReplaceAll(text, "$labels", ".")

	// render
	tmpl, err := template.New("").Parse(text)
	if err != nil {
		return input, fmt.Errorf("error on Parse: %w", err)
	}
	labels := map[string]string{}
	for k, v := range sample.Metric {
		labels[string(k)] = string(v)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, labels)
	if err != nil {
		// test not reachable: buffer full? writer malfunction?
		return input, fmt.Errorf("error on Execute: %w", err)
	}
	return buf.String(), nil
}

func (a *alerter) sendFires(fires []model.Fire) error {
	pbytes, err := json.Marshal(fires)
	if err != nil {
		// test not reachable: memory full?
		return fmt.Errorf("error on Marshal: %w", err)
	}
	buff := bytes.NewBuffer(pbytes)
	resp, err := http.Post(a.alertmanagerURL+"/api/v1/alerts", "application/json", buff)
	if err != nil {
		return fmt.Errorf("error on Post: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("statusCode is not ok(200)")
	}
	return nil
}

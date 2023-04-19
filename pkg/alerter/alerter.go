package alerter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/store"
	commonModel "github.com/prometheus/common/model"
	promRule "github.com/prometheus/prometheus/rules"
)

type alerter struct {
	alertFiles      []model.AlertFile
	remoteStore     *store.RemoteStore
	repeat          bool
	alertmanagerURL string
}

func NewAlerter(stores *store.Stores) *alerter {
	return &alerter{
		alertFiles:      getAlertFiles(stores),
		remoteStore:     stores.RemoteStore,
		alertmanagerURL: "http://alertmanager:9093",
	}
}

func getAlertFiles(stores *store.Stores) []model.AlertFile {
	alertFiles := []model.AlertFile{}
	for _, ruleFile := range stores.AlertRuleStore.AlertRuleFiles() {
		datasources := stores.DatasourceStore.GetDatasourcesWithSelector(ruleFile.DatasourceSelector)
		if len(datasources) < 1 {
			log.Printf("no datasources from GetDatasourcesWithSelector")
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
	log.Printf("starting alerter...")
	a.repeat = true
	go a.loop()
}

func (a *alerter) Stop() {
	log.Printf("stopping alerter...")
	a.repeat = false
}

func (a *alerter) loop() {
	for {
		if !a.repeat {
			log.Printf("alerter stopped")
			return
		}
		a.Once()
		time.Sleep(20 * time.Second)
	}
}

func (a *alerter) Once() {
	log.Printf("processAlertFiles")
	err := a.processAlertFiles()
	if err != nil {
		log.Printf("error on processAlertFiles: %s", err)
	}
}

func (a *alerter) processAlertFiles() error {
	totalFires := []model.Fire{}
	for i := range a.alertFiles {
		for j := range a.alertFiles[i].AlertGroups {
			for k := range a.alertFiles[i].AlertGroups[j].Alerts {
				fires, err := a.processAlert(&a.alertFiles[i].AlertGroups[j].Alerts[k], &a.alertFiles[i].Datasource)
				if err != nil {
					log.Printf("error on processAlert: %s", err)
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
	fires, err := evaluateAlert(alert, queryData)
	if err != nil {
		return zero, fmt.Errorf("error on evaluateAlert: %s", err)
	}
	return fires, nil
}

func (a *alerter) queryAlert(alert *model.Alert, datasource *model.Datasource) (model.QueryData, error) {
	var zero model.QueryData
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	resultString, err := a.remoteStore.Get(ctx, *datasource, "query", "query="+url.QueryEscape(alert.Expr))
	if err != nil {
		return zero, fmt.Errorf("error on remoteStore.Get: %w", err)
	}
	var queryResult model.QueryResult
	//  resultString: {"status":"success","data":{"resultType":"vector","result":[]}}
	err = json.Unmarshal([]byte(resultString), &queryResult)
	if err != nil {
		return zero, fmt.Errorf("error on Unmarshal: %w", err)
	}
	if queryResult.Status != "success" {
		return zero, fmt.Errorf("query status is not success")
	}
	return queryResult.Data, nil
}

func evaluateAlert(alert *model.Alert, queryData model.QueryData) ([]model.Fire, error) {
	var zero []model.Fire

	// inactive
	if len(queryData.Result) < 1 {
		alert.State = promRule.StateInactive
		alert.ActiveAt = 0
		return zero, nil
	}

	if alert.ActiveAt == 0 {
		// now active
		alert.ActiveAt = commonModel.Now()
	}
	// pending
	if alert.ActiveAt.Add(time.Duration(alert.For)).After(commonModel.Now()) {
		alert.State = promRule.StatePending
		return zero, nil
	}
	// firing
	alert.State = promRule.StateFiring
	return getFires(alert, queryData), nil
}

func getFires(alert *model.Alert, data model.QueryData) []model.Fire {
	if len(alert.Labels) < 1 {
		alert.Labels = map[string]string{}
	}
	if len(alert.Annotations) < 1 {
		alert.Annotations = map[string]string{}
	}
	alert.Labels["alertname"] = alert.Name
	alert.Labels["firer"] = "venti"
	if _, exists := alert.Annotations["summary"]; !exists {
		if len(alert.Annotations) == 0 {
			alert.Annotations = map[string]string{}
		}
		alert.Annotations["summary"] = "dummy summary from venti"
	}
	if data.ResultType != commonModel.ValVector {
		log.Printf("resultType is not vector")
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
		fire.Annotations["summary"] = renderSummary(alert.Annotations["summary"], &sample)
		fires = append(fires, fire)
	}
	return fires
}

func renderSummary(tmplString string, sample *commonModel.Sample) string {
	result := tmplString
	result = strings.ReplaceAll(result, "$value", sample.Value.String())
	result = strings.ReplaceAll(result, "$labels.", ".")
	result = strings.ReplaceAll(result, "$labels", ".")
	var buf bytes.Buffer
	t, err := template.New("t1").Parse(result)
	if err != nil {
		log.Printf("error on Parse: %s", err)
		return result
	}

	labels := map[string]string{}
	for k, v := range sample.Metric {
		labels[string(k)] = string(v)
	}
	err = t.Execute(&buf, labels)

	if err != nil {
		log.Printf("error on Execute: %s", err)
		return result
	}
	return buf.String()
}

func (a *alerter) sendFires(fires []model.Fire) error {
	pbytes, err := json.Marshal(fires)
	if err != nil {
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

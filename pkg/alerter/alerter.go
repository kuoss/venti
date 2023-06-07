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
	"github.com/kuoss/venti/pkg/service/alerting"
	"github.com/kuoss/venti/pkg/service/remote"
	commonModel "github.com/prometheus/common/model"
	promRule "github.com/prometheus/prometheus/rules"
)

type alerter struct {
	alertingService      *alerting.AlertingService
	remoteService        *remote.RemoteService
	evaluationInterval time.Duration
	repeat             bool
	alertmanagerURL    string
}

func New(alertingService *alerting.AlertingService, remoteService *remote.RemoteService) *alerter {
	return &alerter{
		alertingService:      alertingService,
		remoteService:        remoteService,
		evaluationInterval: 20 * time.Second,                   // TODO: configurable
		alertmanagerURL:    alertingService.GetAlertmanagerURL(), // TODO: multiple alertmanagers
	}
}

func (a *alerter) SetAlertmanagerURL(url string) {
	a.alertmanagerURL = url
}

func (a *alerter) Start() error {
	if len(a.alertingService.AlertFiles) < 1 {
		return fmt.Errorf("no alertFiles")
	}
	if a.alertmanagerURL == "" {
		logger.Warnf("alertmanagerURL is not set")
	}
	logger.Infof("starting alerter...")
	a.repeat = true
	go a.loop()
	return nil
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
	if a.alertingService == nil {
		return fmt.Errorf("nil alertingService")
	}
	if len(a.alertingService.AlertFiles) < 1 {
		return fmt.Errorf("no alert files")
	}
	var fires []model.Fire
	for i := range a.alertingService.AlertFiles {
		for j := range a.alertingService.AlertFiles[i].AlertGroups {
			for k := range a.alertingService.AlertFiles[i].AlertGroups[j].RuleAlerts {
				temps := a.processRuleAlert(&a.alertingService.AlertFiles[i].AlertGroups[j].RuleAlerts[k], &a.alertingService.AlertFiles[i].CommonLabels)
				fires = append(fires, temps...)
			}
		}
	}
	err := a.sendFires(fires)
	if err != nil {
		return fmt.Errorf("sendFires err: %w", err)
	}
	return nil
}

func (a *alerter) processRuleAlert(ruleAlert *model.RuleAlert, commonLabels *map[string]string) []model.Fire {
	var fires []model.Fire
	rule := &ruleAlert.Rule
	for i := range ruleAlert.Alerts {
		alert := &ruleAlert.Alerts[i]
		queryData, err := a.queryAlert(rule, alert)
		if err != nil {
			logger.Warnf("queryAlert err: %s", err.Error())
			continue
		}
		temps := evaluateAlert(queryData, rule, alert, commonLabels)
		fires = append(fires, temps...)
	}
	return fires
}

func (a *alerter) queryAlert(rule *model.Rule, alert *model.Alert) (model.QueryData, error) {
	if alert.Datasource == nil {
		return model.QueryData{}, fmt.Errorf("datasource is nil")
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	code, body, err := a.remoteService.GET(ctx, alert.Datasource, remote.ActionQuery, "query="+url.QueryEscape(rule.Expr))
	if err != nil {
		return model.QueryData{}, fmt.Errorf("GET err: %w", err)
	}

	//  body: {"status":"success","data":{"resultType":"vector","result":[]}}
	var queryResult model.QueryResult
	err = json.Unmarshal([]byte(body), &queryResult)
	if err != nil {
		return model.QueryData{}, fmt.Errorf("unmarshal err: %s, body: %s", err.Error(), body)
	}
	// wrap up
	if queryResult.Status != "success" {
		err = fmt.Errorf("not success status (status=%s, code=%d)", queryResult.Status, code)
	} else {
		if code != http.StatusOK {
			// test not reachable: non-200 {"status":"success"}
			err = fmt.Errorf("not ok (code=%d)", code)
		}
	}
	return queryResult.Data, err
}

func evaluateAlert(queryData model.QueryData, rule *model.Rule, alert *model.Alert, commonLabels *map[string]string) []model.Fire {
	var fires []model.Fire
	// inactive
	if len(queryData.Result) < 1 {
		alert.State = promRule.StateInactive
		alert.ActiveAt = 0
		return fires
	}
	if alert.ActiveAt == 0 {
		// now active
		alert.ActiveAt = commonModel.Now()
	}
	// pending
	if alert.ActiveAt.Add(time.Duration(rule.For)).After(commonModel.Now()) {
		alert.State = promRule.StatePending
		logger.Infof("evaluateAlert: [pending] %s", rule.Alert)
		return fires
	}
	// firing
	logger.Infof("evaluateAlert: [firing] %s", rule.Alert)
	alert.State = promRule.StateFiring
	return getFires(rule, queryData, commonLabels)
}

func getFires(rule *model.Rule, data model.QueryData, commonLabels *map[string]string) []model.Fire {
	labels := map[string]string{}
	annotations := map[string]string{}
	if commonLabels != nil {
		for k, v := range *commonLabels {
			labels[k] = v
		}
	}
	if rule.Labels != nil {
		for k, v := range rule.Labels {
			labels[k] = v
		}
	}
	if rule.Annotations != nil {
		for k, v := range rule.Annotations {
			annotations[k] = v
		}
	}
	labels["alertname"] = rule.Alert
	if labels["alertname"] == "" {
		labels["alertname"] = "placeholder name"
	}
	labels["firer"] = "venti"

	if _, exists := annotations["summary"]; !exists {
		annotations["summary"] = "placeholder summary"
	}
	if data.ResultType != commonModel.ValVector {
		logger.Warnf("resultType is not vector")
		return []model.Fire{{
			Labels:      labels,
			Annotations: annotations,
		}}
	}

	fires := []model.Fire{}
	for _, sample := range data.Result {
		fire := model.Fire{
			Labels:      map[string]string{},
			Annotations: map[string]string{},
		}
		// deep copy
		for k, v := range labels {
			fire.Labels[k] = v
		}
		for k, v := range annotations {
			fire.Annotations[k] = v
		}
		// render summary
		summary, err := renderSummary(annotations["summary"], &sample)
		if err != nil {
			logger.Warnf("renderSummary err: %s", err.Error())
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

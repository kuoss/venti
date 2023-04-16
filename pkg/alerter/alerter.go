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
	"reflect"
	"strings"
	"time"

	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/store"
	commonModel "github.com/prometheus/common/model"
	promRule "github.com/prometheus/prometheus/rules"
)

type alerter struct {
	alertFiles  []model.AlertFile
	remoteStore *store.RemoteStore
	repeat      bool
}

func NewAlerter(stores *store.Stores) *alerter {
	return &alerter{
		alertFiles:  getAlertFiles(stores),
		remoteStore: stores.RemoteStore,
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
					alerts = append(alerts, model.Alert{
						State:       promRule.StateInactive,
						Name:        rule.Alert,
						Expr:        rule.Expr,
						For:         rule.For,
						Labels:      rule.Labels,
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

func (a *alerter) Start() {
	log.Println("starting alerter...")
	a.repeat = true
	go a.run()
}

func (a *alerter) Stop() {
	log.Println("stopping alerter...")
	a.repeat = false
}

func (a *alerter) run() {
	for {
		if !a.repeat {
			log.Println("alerter stopped")
			return
		}
		log.Println("alerter task")
		a.Once()
		time.Sleep(20 * time.Second)
	}
}

func (a *alerter) Once() {
	a.evaluateAlertFiles()
	a.fireAlertFiles()
}

func (a *alerter) evaluateAlertFiles() {
	for i := range a.alertFiles {
		for j := range a.alertFiles[i].AlertGroups {
			for k := range a.alertFiles[i].AlertGroups[j].Alerts {
				err := a.evaluateAlert(&a.alertFiles[i].AlertGroups[j].Alerts[k], a.alertFiles[i].Datasource)
				if err != nil {
					log.Fatalf("error on evaluateAlert: %s", err.Error())
				}
				time.Sleep(time.Duration(500) * time.Millisecond)
			}
		}
	}
}

func (a *alerter) evaluateAlert(alert *model.Alert, datasource model.Datasource) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resultString, err := a.remoteStore.Get(ctx, datasource, "query", "query="+url.QueryEscape(alert.Expr))
	if err != nil {
		return fmt.Errorf("error on remoteStore.Get: %w", err)
	}
	var queryResult model.QueryResult
	//  resultString: {"status":"success","data":{"resultType":"vector","result":[]}}
	err = json.Unmarshal([]byte(resultString), &queryResult)
	if err != nil {
		return fmt.Errorf("error on Unmarshal: %w", err)
	}
	if queryResult.Status != "success" {
		return fmt.Errorf("query status is not success")
	}

	// inactive
	if len(queryResult.Data.Result) < 1 {
		alert.State = promRule.StateInactive
		alert.ActiveAt = 0
		return nil
	}

	// set ActiveAt if it is zero
	if reflect.ValueOf(alert.ActiveAt).IsZero() {
		alert.ActiveAt = commonModel.Now()
	}
	// pending
	if alert.ActiveAt.Add(time.Duration(alert.For)).After(commonModel.Now()) {
		alert.State = promRule.StatePending
		return nil
	}
	// firing: add to firing alerts
	alert.QueryData = queryResult.Data
	alert.State = promRule.StateFiring
	return nil
}

func (a *alerter) fireAlertFiles() {
	fires := []model.Fire{}
	for _, alertFile := range a.alertFiles {
		for _, group := range alertFile.AlertGroups {
			for _, alert := range group.Alerts {
				if alert.State == promRule.StateFiring {
					fires = append(fires, a.getFiresFromAlert(alert)...)
				}
			}
		}
	}
	cnt := len(fires)
	if cnt < 1 {
		return
	}
	err := a.fireFires(fires)
	if err != nil {
		log.Fatalf("error on fireFires: %s", err.Error())
		return
	}
	log.Printf("fireFires success: %d fires\n", cnt)
}

func (a *alerter) getFiresFromAlert(alert model.Alert) []model.Fire {
	alert.Labels["alertname"] = alert.Name
	alert.Labels["firer"] = "venti"
	if _, exists := alert.Annotations["summary"]; !exists {
		alert.Annotations["summary"] = "dummy summary from venti"
	}
	if alert.QueryData.ResultType != commonModel.ValVector {
		log.Println("resultType is not vector")
		return []model.Fire{{
			State:       "firing",
			Labels:      alert.Labels,
			Annotations: alert.Annotations,
		}}
	}

	fires := []model.Fire{}
	for _, sample := range alert.QueryData.Result {
		// deep copy alert
		fire := model.Fire{State: "firing", Labels: map[string]string{}, Annotations: map[string]string{}}
		for k, v := range alert.Labels {
			fire.Labels[k] = v
		}
		for k, v := range alert.Annotations {
			fire.Annotations[k] = v
		}
		fire.Annotations["summary"] = a.renderSummary(alert.Annotations["summary"], &sample)
		fires = append(fires, fire)
	}
	return fires
}

func (a *alerter) renderSummary(tmplString string, sample *commonModel.Sample) string {
	result := tmplString
	result = strings.ReplaceAll(result, "$value", sample.Value.String())
	result = strings.ReplaceAll(result, "$labels.", ".")
	result = strings.ReplaceAll(result, "$labels", ".")
	var buf bytes.Buffer
	t, err := template.New("t1").Parse(result)
	if err != nil {
		log.Printf("warn: cannot render: %s", tmplString)
		return result
	}
	err = t.Execute(&buf, sample.Metric)
	if err != nil {
		log.Printf("warn: cannot render: %s", tmplString)
		return result
	}
	return buf.String()
}

func (a *alerter) fireFires(shots []model.Fire) error {
	pbytes, err := json.Marshal(shots)
	if err != nil {
		return fmt.Errorf("error on Marshal: %w", err)
	}
	buff := bytes.NewBuffer(pbytes)
	resp, err := http.Post("http://alertmanager:9093/api/v1/alerts", "application/json", buff)
	if err != nil {
		return fmt.Errorf("error on Post: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("statusCode is not ok(200)")
	}
	return nil
}

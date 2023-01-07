package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type ValueType string

const (
	ValueTypeNone   ValueType = "none"
	ValueTypeVector ValueType = "vector"
	ValueTypeScalar ValueType = "scalar"
	ValueTypeMatrix ValueType = "matrix"
	ValueTypeString ValueType = "string"

	ValueTypeLogs ValueType = "logs"
)

type Sample struct {
	Vaue   []interface{}     `json:"value"`
	Metric map[string]string `json:"metric"`
}

type Vector []Sample

// {"data":{"result":[],"resultType":"logs"},"status":"success"}
// {"status":"success","data":{"resultType":"vector","result":[]}}
type QueryResult struct {
	Data   QueryData `json:"data"`
	Status string    `json:"status"`
}
type QueryData struct {
	ResultType ValueType `json:"resultType"`
	Result     Vector    `json:"result"`
}

type Alert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels,omitempty"`
	Annotations  map[string]string `json:"annotations,omitempty"`
	GeneratorURL string            `json:"generatorURL,omitempty"`
}

func StartAlertDaemon() {
	log.Println("starting alert daemon...")
	go alertTaskLoop()
	go alertLogTaskLoop()
}

/*
//todo refactoring
				log.Println("----Workaround for log alert")
				logAlert := Alert{Status: "firing", Labels: map[string]string{}, Annotations: map[string]string{}}
				logAlert.Labels = rule.Labels

				// result ->  marshal  -> unmarshal to LogResultType

				var logResult logResult
				err = json.Unmarshal([]byte(resultString), &logResult)
				if err != nil {
					log.Println(err)
				}
				logAlert.Annotations["summary"] = logResult.Data.Result[len(logResult.Data.Result)-1]
				logAlert.Status = "firing"
				logAlert.Labels["alertname"] = rule.Alert
				logAlert.Labels["firer"] = "venti"
				firingAlerts = append(firingAlerts, logAlert)
*/

func alertLogTaskLoop() {
	for {
		log.Println("alert Log task at", time.Now())
		alertLogTask()
		time.Sleep(time.Duration(30) * time.Second)
	}
}

func alertLogTask() {
	alertRuleGroups := GetAlertRuleGroups()
	firingAlerts := []Alert{}
	for groupIndex, group := range alertRuleGroups {
		for ruleIndex, rule := range group.Rules {
			now := time.Now()
			if group.DatasourceType != DatasourceTypeLethe {
				continue
			}
			time.Sleep(time.Duration(1500) * time.Millisecond)

			resultString, err := RunHTTPLetheQuery(HTTPQuery{
				Query: rule.Expr,
			})

			var logResult logResult
			err = json.Unmarshal([]byte(resultString), &logResult)
			if err != nil {
				log.Println(err)
				logResult.Data = logData{
					ResultType: ValueTypeLogs,
					Result:     []string{},
				}
			}

			if len(logResult.Data.Result) < 1 {
				fmt.Println("result empty. inactivating")
				// NORMAL
				config.AlertRuleGroups[groupIndex].Rules[ruleIndex].State = AlertStateInactive
				config.AlertRuleGroups[groupIndex].Rules[ruleIndex].ActiveAt = time.Time{}
				continue
			}
			if reflect.ValueOf(rule.ActiveAt).IsZero() {
				fmt.Println("result empty. inactivating")
				config.AlertRuleGroups[groupIndex].Rules[ruleIndex].ActiveAt = now
			}
			// pending
			if config.AlertRuleGroups[groupIndex].Rules[ruleIndex].ActiveAt.Add(rule.For).After(now) {
				fmt.Println("make rule pending")
				config.AlertRuleGroups[groupIndex].Rules[ruleIndex].State = AlertStatePending
				continue
			}
			// firing: add to firing alerts
			config.AlertRuleGroups[groupIndex].Rules[ruleIndex].State = AlertStateFiring

			logAlert := Alert{
				Status:      "firing",
				Labels:      rule.Labels,
				Annotations: map[string]string{},
			}

			logAlert.Labels["alertname"] = rule.Alert
			logAlert.Labels["firer"] = "venti"
			logAlert.Annotations["summary"] = logResult.Data.Result[len(logResult.Data.Result)-1]
			firingAlerts = append(firingAlerts, logAlert)
			fmt.Printf("logAlert added %+vln", logAlert)
		}
	}
	fireAlerts(firingAlerts)
}

func alertTaskLoop() {
	for {
		log.Println("alert task at", time.Now())
		alertTask()
		time.Sleep(time.Duration(30) * time.Second)
	}
}

func alertTask() {
	var resultString string
	var err error
	alertRuleGroups := GetAlertRuleGroups()
	firingAlerts := []Alert{}
	for i, group := range alertRuleGroups {
		datasourceType := group.DatasourceType
		for j, rule := range group.Rules {
			time.Sleep(time.Duration(500) * time.Millisecond)
			now := time.Now()
			switch datasourceType {
			case DatasourceTypeLethe:
				resultString, err = RunHTTPLetheQuery(HTTPQuery{
					Query: rule.Expr,
				})
				continue

			case DatasourceTypePrometheus:
				resultString, err = RunHTTPPrometheusQuery(HTTPQuery{
					Query: rule.Expr,
				})
			}
			if err != nil {
				log.Printf("error on query: err=%s, expr=%s\n", err, rule.Expr)
				continue
			}
			// log.Println(resultString)
			var queryResult QueryResult
			err = json.Unmarshal([]byte(resultString), &queryResult)
			if err != nil {
				log.Printf("unmarshal error: err=%s, expr=%s\n", err, rule.Expr)
				continue
			}
			if queryResult.Status != "success" {
				log.Println("query status is not success")
				continue
			}
			if len(queryResult.Data.Result) < 1 {
				// NORMAL
				config.AlertRuleGroups[i].Rules[j].State = AlertStateInactive
				config.AlertRuleGroups[i].Rules[j].ActiveAt = time.Time{}
				continue
			}
			if reflect.ValueOf(rule.ActiveAt).IsZero() {
				config.AlertRuleGroups[i].Rules[j].ActiveAt = now
			}
			// pending
			if config.AlertRuleGroups[i].Rules[j].ActiveAt.Add(rule.For).After(now) {
				config.AlertRuleGroups[i].Rules[j].State = AlertStatePending
				continue
			}
			// firing: add to firing alerts
			config.AlertRuleGroups[i].Rules[j].State = AlertStateFiring

			firingAlerts = append(firingAlerts, renderAlerts(rule, queryResult)...)

		}
	}
	fireAlerts(firingAlerts)
}

func renderTemplate(tmplString string, labels map[string]string, value string) string {
	result := tmplString
	result = strings.ReplaceAll(result, "$value", value)
	result = strings.ReplaceAll(result, "$labels.", ".")
	result = strings.ReplaceAll(result, "$labels", ".")
	var buf bytes.Buffer
	t, err := template.New("t1").Parse(result)
	if err != nil {
		log.Printf("warn: cannot render: %s", tmplString)
		return result
	}
	err = t.Execute(&buf, labels)
	if err != nil {
		log.Printf("warn: cannot render: %s", tmplString)
		return result
	}
	return buf.String()
}

type logResult struct {
	Data   logData `json:"data"`
	Status string  `json:"status"`
}
type logData struct {
	ResultType ValueType `json:"resultType"`
	Result     []string  `json:"result"`
}

func renderAlerts(rule AlertRule, result QueryResult) []Alert {
	rule.Labels["alertname"] = rule.Alert
	rule.Labels["firer"] = "venti"

	if _, exists := rule.Annotations["summary"]; !exists {
		rule.Annotations["summary"] = "venti alerting.go"
	}

	if result.Data.ResultType != ValueTypeVector {
		log.Println("ResultType is not Vector ㅠㅠ")
		return []Alert{{
			Status:      "firing",
			Labels:      rule.Labels,
			Annotations: rule.Annotations,
		}}
	}

	alerts := []Alert{}
	for _, smpl := range result.Data.Result {
		// deep copy alert
		alert := Alert{Status: "firing", Labels: map[string]string{}, Annotations: map[string]string{}}
		for k, v := range rule.Labels {
			alert.Labels[k] = v
		}
		for k, v := range rule.Annotations {
			alert.Annotations[k] = v
		}
		alert.Annotations["summary"] = renderTemplate(alert.Annotations["summary"], smpl.Metric, smpl.Vaue[1].(string))
		// log.Println(alert.Annotations["summary"])
		alerts = append(alerts, alert)
	}
	return alerts
}

func fireAlerts(alerts []Alert) {
	if len(alerts) < 1 {
		return
	}
	pbytes, err := json.Marshal(alerts)
	if err != nil {
		log.Println("alerting.go: cannot marshal alerts")
		return
	}
	buff := bytes.NewBuffer(pbytes)
	response, err := http.Post("http://alertmanager:9093/api/v1/alerts", "application/json", buff)
	if err != nil {
		log.Printf("alertmanager failed: %s", err)
		return
	}
	log.Println("alertmanager", response.StatusCode)
}

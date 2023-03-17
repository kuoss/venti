package server

import (
	"bytes"
	"encoding/json"
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
}

func alertTaskLoop() {
	for {
		log.Println("alert task at", time.Now())
		alertTask()
		time.Sleep(time.Duration(30) * time.Second)
	}
}

func alertTask() {
	alertRuleGroups := GetAlertRuleGroups()
	firingAlerts := []Alert{}
	for i, group := range alertRuleGroups {
		datasourceType := group.DatasourceType
		for j, rule := range group.Rules {
			time.Sleep(time.Duration(500) * time.Millisecond)
			now := time.Now()
			instanctQuery := InstantQuery{
				DatasourceType: datasourceType,
				Expr: rule.Expr,
			}
			resultString, err := instanctQuery.execute()
			if err != nil {
				log.Printf("error on query: err=%s, expr=%s\n", err, rule.Expr)
				continue
			}
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

package server

import (
	"bytes"
	"encoding/json"
	"errors"
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

func renderTemplate(tmplString string, labels map[string]string) (string, error) {
	tmplString = strings.ReplaceAll(tmplString, "$labels.", ".")
	var buf bytes.Buffer
	t, err := template.New("t1").Parse(tmplString)
	if err != nil {
		return "", errors.New("template parse error")
	}
	err = t.Execute(&buf, labels)
	if err != nil {
		return "", errors.New("template execute error")
	}
	return buf.String(), nil
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
	var err error
	for _, smpl := range result.Data.Result {
		// deep copy alert
		alert := Alert{Status: "firing", Labels: map[string]string{}, Annotations: map[string]string{}}
		for k, v := range rule.Labels {
			alert.Labels[k] = v
		}
		for k, v := range rule.Annotations {
			alert.Annotations[k] = v
		}
		alert.Annotations["summary"], err = renderTemplate(alert.Annotations["summary"], smpl.Metric)
		if err != nil {
			log.Printf("cannot render summary: %s", err)
			continue
		}
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

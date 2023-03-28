package alert

import (
	"bytes"
	"encoding/json"
	"github.com/kuoss/venti/server"
	"github.com/kuoss/venti/server/configuration"
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

type AlertRuleGroupList struct {
	Groups []AlertRuleGroup `json:"groups"`
}

type AlertRuleGroup struct {
	Name           string                `json:"name"`
	Rules          []AlertRule           `json:"rules"`
	DatasourceType server.DatasourceType `json:"datasource" yaml:"datasource"`
	CommonLabels   map[string]string     `json:"commonLabels,omitempty" yaml:"commonLabels,omitempty"`
}

type AlertRule struct {
	Alert       string            `json:"alert,omitempty"`
	Expr        string            `json:"expr"`
	For         time.Duration     `json:"for,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	State       AlertState        `json:"state,omitempty"`
	ActiveAt    time.Time         `json:"activeStartTime,omitempty"`
}

type AlertState string

const (
	AlertStateInactive AlertState = "inactive"
	AlertStatePending  AlertState = "pending"
	AlertStateFiring   AlertState = "firing"
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
	alertRuleGroups := configuration.GetAlertRuleGroups()
	firingAlerts := []Alert{}
	for i, group := range alertRuleGroups {
		datasourceType := group.DatasourceType
		for j, rule := range group.Rules {
			time.Sleep(time.Duration(500) * time.Millisecond)
			now := time.Now()
			instanctQuery := server.InstantQuery{
				DatasourceType: datasourceType,
				Expr:           rule.Expr,
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
				configuration.config.AlertRuleGroups[i].Rules[j].State = AlertStateInactive
				configuration.config.AlertRuleGroups[i].Rules[j].ActiveAt = time.Time{}
				continue
			}
			if reflect.ValueOf(rule.ActiveAt).IsZero() {
				configuration.config.AlertRuleGroups[i].Rules[j].ActiveAt = now
			}
			// pending
			if configuration.config.AlertRuleGroups[i].Rules[j].ActiveAt.Add(rule.For).After(now) {
				configuration.config.AlertRuleGroups[i].Rules[j].State = AlertStatePending
				continue
			}
			// firing: add to firing alerts
			configuration.config.AlertRuleGroups[i].Rules[j].State = AlertStateFiring
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
		rule.Annotations["summary"] = "venti alert.go"
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
		log.Println("alert.go: cannot marshal alerts")
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

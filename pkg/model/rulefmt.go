package model

import (
	commonModel "github.com/prometheus/common/model"
	promRule "github.com/prometheus/prometheus/rules"
)

/*
types from prometheus projects
- Let's consider the commonModel first! 🚀

https://github.com/prometheus/alertmanager/blob/v0.25.0/api/v2/models/alert.go         | Alert
https://github.com/prometheus/alertmanager/blob/v0.25.0/api/v2/models/alert_group.go   | AlertGroup
https://github.com/prometheus/alertmanager/blob/v0.25.0/api/v2/models/alert_groups.go  | AlertGroups
https://github.com/prometheus/alertmanager/blob/v0.25.0/types/types.go                 | AlertState
https://github.com/prometheus/common/blob/v0.42.0/model/value.go                       | 🚀 Sample, Samples, Scalar, String, Vector, Matrix, Value✔️, ValueType✔️, ValueTypeVector✔️
https://github.com/prometheus/prometheus/blob/v2.43.0/model/rulefmt/rulefmt.go         | RuleGroup➡️RuleGroup, Rule➡️Rule
https://github.com/prometheus/prometheus/blob/v2.43.0/notifier/notifier.go             | Alert
https://github.com/prometheus/prometheus/blob/v2.43.0/promql/value.go                  | Sample, Vector
https://github.com/prometheus/prometheus/blob/v2.43.0/promql/parser/value.go           | Value, ValueType, ValueTypeVector
https://github.com/prometheus/prometheus/blob/v2.43.0/rules/alerting.go                | Alert➡️Alert, AlertState✔️, AlertingRule➡️Alert
https://github.com/prometheus/prometheus/blob/v2.43.0/template/template.go             | sample, queryResult
https://github.com/prometheus/prometheus/blob/v2.43.0/web/api/v1/api.go                | apiFuncResult➡️QueryResult, queryData➡️QueryData
*/

type RuleFile struct {
	Kind               string             `json:"kind,omitempty" yaml:"kind,omitempty"`
	CommonLabels       map[string]string  `json:"commonLabels,omitempty" yaml:"commonLabels,omitempty"`
	DatasourceSelector DatasourceSelector `json:"datasourceSelector" yaml:"datasourceSelector"`
	RuleGroups         []RuleGroup        `json:"groups" yaml:"groups"`
}

type AlertFile struct {
	Kind               string             `json:"kind,omitempty"`
	CommonLabels       map[string]string  `json:"commonLabels,omitempty"`
	DatasourceSelector DatasourceSelector `json:"datasourceSelector"`
	AlertGroups        []AlertGroup       `json:"groups"`
}

type AlertGroup struct {
	Name       string               `json:"name"`
	Interval   commonModel.Duration `json:"interval,omitempty"`
	Limit      int                  `json:"limit,omitempty"`
	RuleAlerts []RuleAlert          `json:"ruleAlerts"`
}

type RuleAlert struct {
	Rule   Rule    `json:"rule"`
	Alerts []Alert `json:"alerts"`
}

type Alert struct {
	Datasource *Datasource         `json:"datasource"`
	State      promRule.AlertState `json:"state"` // commonModel doesn't have pending state
	ActiveAt   commonModel.Time    `json:"activeAt"`
}

type Fire struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

type RuleGroup struct {
	Name     string               `json:"name" yaml:"name"`
	Interval commonModel.Duration `json:"interval,omitempty" yaml:"interval,omitempty"`
	Limit    int                  `json:"limit,omitempty" yaml:"limit,omitempty"`
	Rules    []Rule               `json:"rules" yaml:"rules"`
}

type Rule struct {
	Record        string               `json:"record,omitempty" yaml:"record,omitempty"`
	Alert         string               `json:"alert,omitempty" yaml:"alert,omitempty"`
	Expr          string               `json:"expr" yaml:"expr"`
	For           commonModel.Duration `json:"for" yaml:"for,omitempty"`
	KeepFiringFor commonModel.Duration `json:"keep_firing_for,omitempty" yaml:"keep_firing_for,omitempty"`
	Labels        map[string]string    `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations   map[string]string    `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

type QueryResult struct {
	Data   QueryData `json:"data"`
	Status string    `json:"status"`
}

type QueryData struct {
	ResultType commonModel.ValueType `json:"resultType"`
	Result     []commonModel.Sample  `json:"result"`
}

// resultString: {"status":"success","data":{"resultType":"vector","result":[]}}

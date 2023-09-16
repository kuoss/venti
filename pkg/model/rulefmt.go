package model

import (
	"time"
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

// https://github.com/prometheus/prometheus/blob/v2.43.0/model/rulefmt/rulefmt.go#L137
type RuleGroup struct {
	Name     string        `json:"name" yaml:"name"`
	Interval time.Duration `json:"interval,omitempty" yaml:"interval,omitempty"`
	Limit    int           `json:"limit,omitempty" yaml:"limit,omitempty"`
	Rules    []Rule        `json:"rules" yaml:"rules"`
}

type Rule struct {
	Record        string            `json:"record,omitempty" yaml:"record,omitempty"`
	Alert         string            `json:"alert,omitempty" yaml:"alert,omitempty"`
	Expr          string            `json:"expr" yaml:"expr"`
	For           time.Duration     `json:"for" yaml:"for,omitempty"`
	KeepFiringFor time.Duration     `json:"keep_firing_for,omitempty" yaml:"keep_firing_for,omitempty"`
	Labels        map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations   map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

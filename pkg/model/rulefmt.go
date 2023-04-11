package model

import (
	commonModel "github.com/prometheus/common/model"
)

type RuleFile struct {
	Kind               string             `json:"kind,omitempty" yaml:"kind,omitempty"`
	CommonLabels       map[string]string  `json:"commonLabels,omitempty" yaml:"commonLabels,omitempty"`
	DatasourceSelector DatasourceSelector `json:"datasourceSelector" yaml:"datasourceSelector"`
	Groups             []RuleGroup        `json:"groups" yaml:"groups"`
}

// Prometheus rulefmt doesn't have json annotations.
// https://github.com/prometheus/prometheus/blob/main/model/rulefmt/rulefmt.go

// RuleGroup is a list of sequentially evaluated recording and alerting rules.
type RuleGroup struct {
	Name     string               `json:"name" yaml:"name"`
	Interval commonModel.Duration `json:"interval,omitempty" yaml:"interval,omitempty"`
	Limit    int                  `json:"limit,omitempty" yaml:"limit,omitempty"`
	Rules    []Rule               `json:"rules" yaml:"rules"`
}

// Rule describes an alerting or recording rule.
type Rule struct {
	Record        string               `json:"record,omitempty" yaml:"record,omitempty"`
	Alert         string               `json:"alert,omitempty" yaml:"alert,omitempty"`
	Expr          string               `json:"expr" yaml:"expr"`
	For           commonModel.Duration `json:"for,omitempty" yaml:"for,omitempty"`
	KeepFiringFor commonModel.Duration `json:"keep_firing_for,omitempty" yaml:"keep_firing_for,omitempty"`
	Labels        map[string]string    `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations   map[string]string    `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

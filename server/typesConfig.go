package server

import (
	"time"
)

type Config struct {
	Version           string
	EtcUsersConfig    EtcUsersConfig
	DatasourcesConfig DatasourcesConfig
	Dashboards        []Dashboard
	AlertRuleGroups   []AlertRuleGroup
}

// user
type EtcUsersConfig struct {
	EtcUsers []EtcUser `yaml:"users"`
}

type EtcUser struct {
	Username string `yaml:"username"`
	Hash     string `yaml:"hash"`
	IsAdmin  bool   `yaml:"isAdmin,omitempty"`
}

// datasource
type DatasourceType string

const (
	DatasourceTypePrometheus DatasourceType = "Prometheus"
	DatasourceTypeLethe      DatasourceType = "Lethe"
)

type DatasourcesConfig struct {
	QueryTimeout time.Duration `json:"queryTimeout"`
	Datasources  []Datasource  `json:"datasources"`
}

type Datasource struct {
	Type         DatasourceType `json:"type"`
	Host         string         `json:"host"`
	Port         int            `json:"port,omitempty" yaml:"port,omitempty"`
	IsDiscovered bool           `json:"is_discovered,omitempty" yaml:"isDiscovered,omitempty"`
}

// dashboard
type Dashboard struct {
	Title string `json:"title"`
	Rows  []Row  `json:"rows"`
}

type Row struct {
	Panels []Panel `json:"panels"`
}

type Panel struct {
	Title        string        `json:"title" yaml:"title"`
	Type         string        `json:"type" yaml:"type"`
	Headers      []string      `json:"headers,omitempty" yaml:"headers,omitempty"`
	Targets      []Target      `json:"targets" yaml:"targets"`
	ChartOptions *ChartOptions `json:"chartOptions,omitempty" yaml:"chartOptions,omitempty"`
}

type ChartOptions struct {
	YMax int `json:"yMax,omitempty" yaml:"yMax,omitempty"`
}

type Target struct {
	Expr       string      `json:"expr"`
	Legend     string      `json:"legend,omitempty" yaml:"legend,omitempty"`
	Legends    []string    `json:"legends,omitempty" yaml:"legends,omitempty"`
	Unit       string      `json:"unit,omitempty" yaml:"unit,omitempty"`
	Columns    []string    `json:"columns,omitempty" yaml:"columns,omitempty"`
	Headers    []string    `json:"headers,omitempty" yaml:"headers,omitempty"`
	Key        string      `json:"key,omitempty" yaml:"key,omitempty"`
	Thresholds []Threshold `json:"thresholds,omitempty" yaml:"thresholds,omitempty"`
}

type Threshold struct {
	Values []int `yaml:"values,omitempty" json:"values,omitempty"`
	Invert bool  `yaml:"invert,omitempty" json:"invert,omitempty"`
}

type AlertRuleGroupList struct {
	Groups []AlertRuleGroup `json:"groups"`
}

type AlertRuleGroup struct {
	Name           string            `json:"name"`
	Rules          []AlertRule       `json:"rules"`
	DatasourceType DatasourceType    `json:"datasource" yaml:"datasource"`
	CommonLabels   map[string]string `json:"commonLabels,omitempty" yaml:"commonLabels,omitempty"`
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

package model

import (
	"time"

	commonModel "github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/model/relabel"
)

type Config struct {
	Version          string
	GlobalConfig     GlobalConfig
	DatasourceConfig DatasourceConfig
	UserConfig       UserConfig
	AlertingConfig   AlertingConfig
}

type GlobalConfig struct {
	EvaluationInterval time.Duration `yaml:"evaluation_interval,omitempty"`
}

type UserConfig struct {
	EtcUsers []EtcUser `yaml:"users"`
}

type EtcUser struct {
	Username string `yaml:"username"`
	Hash     string `yaml:"hash"`
	IsAdmin  bool   `yaml:"isAdmin,omitempty"`
}

type DatasourceConfig struct {
	QueryTimeout time.Duration `json:"queryTimeout,omitempty" yaml:"queryTimeout,omitempty"`
	Datasources  []Datasource  `json:"datasources" yaml:"datasources,omitempty"`
	Discovery    Discovery     `json:"discovery,omitempty" yaml:"discovery,omitempty"`
}

type Discovery struct {
	Enabled          bool   `json:"enabled,omitempty" yaml:"enabled,omitempty"`                   // default: false
	MainNamespace    string `json:"mainNamespace,omitempty" yaml:"mainNamespace,omitempty"`       // default: ''
	AnnotationKey    string `json:"annotationKey,omitempty" yaml:"annotationKey,omitempty"`       // default: kuoss.org/datasource-type
	ByNamePrometheus bool   `json:"byNamePrometheus,omitempty" yaml:"byNamePrometheus,omitempty"` // deprecated
	ByNameLethe      bool   `json:"byNameLethe,omitempty" yaml:"byNameLethe,omitempty"`           // deprecated
}

// AlertingConfig...
// https://github.com/prometheus/prometheus/blob/main/config/config.go

type AlertingConfigFile struct {
	AlertingConfig AlertingConfig `yaml:"alerting,omitempty"`
}

type AlertingConfig struct {
	AlertRelabelConfigs []*relabel.Config   `yaml:"alert_relabel_configs,omitempty"`
	AlertmanagerConfigs AlertmanagerConfigs `yaml:"alertmanagers,omitempty"`
}

// AlertmanagerConfigs is a slice of *AlertmanagerConfig.
type AlertmanagerConfigs []*AlertmanagerConfig
type AlertmanagerConfig struct {
	StaticConfig []*TargetGroup `yaml:"static_configs,omitempty"`
}
type TargetGroup struct {
	Targets []string             `yaml:"targets"`
	Labels  commonModel.LabelSet `yaml:"labels,omitempty"`
}

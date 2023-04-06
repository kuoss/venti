package model

import "time"

type Config struct {
	Version           string
	UserConfig        UsersConfig
	DatasourcesConfig *DatasourcesConfig
}

type UsersConfig struct {
	EtcUsers []EtcUser `yaml:"users"`
}

type EtcUser struct {
	Username string `yaml:"username"`
	Hash     string `yaml:"hash"`
	IsAdmin  bool   `yaml:"isAdmin,omitempty"`
}

type DatasourcesConfig struct {
	QueryTimeout time.Duration `json:"queryTimeout,omitempty" yaml:"queryTimeout,omitempty"`
	Datasources  []*Datasource `json:"datasources" yaml:"datasources,omitempty"`
	Discovery    Discovery     `json:"discovery,omitempty" yaml:"discovery,omitempty"`
}

type Discovery struct {
	Enabled          bool   `json:"enabled,omitempty" yaml:"enabled,omitempty"`                   // default: false
	MainNamespace    string `json:"mainNamespace,omitempty" yaml:"mainNamespace,omitempty"`       // default: ''
	AnnotationKey    string `json:"annotationKey,omitempty" yaml:"annotationKey,omitempty"`       // default: kuoss.org/datasource-type
	ByNamePrometheus bool   `json:"byNamePrometheus,omitempty" yaml:"byNamePrometheus,omitempty"` // deprecated
	ByNameLethe      bool   `json:"byNameLethe,omitempty" yaml:"byNameLethe,omitempty"`           // deprecated
}

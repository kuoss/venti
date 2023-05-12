package model

type AlertingFile struct {
	Alertings []Alerting `yaml:"alertings"`
}

type Alerting struct {
	Name string       `yaml:"name"`
	Type AlertingType `yaml:"type"`
	URL  string       `yaml:"url"`
}

type AlertingType string

const (
	AlertingTypeAlertmanager AlertingType = "alertmanager"
)

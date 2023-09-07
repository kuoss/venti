package alerting

import (
	"time"

	"github.com/kuoss/venti/pkg/model"
)

// https://github.com/prometheus/prometheus/blob/v2.46.0/rules/alerting.go

type AlertState int

const (
	StateInactive AlertState = iota
	StatePending
	StateFiring
)

type Alert struct {
	State     AlertState
	CreatedAt time.Time
	UpdatedAt time.Time

	Labels      map[string]string
	Annotations map[string]string
}

type AlertingRule struct {
	datasourceSelector model.DatasourceSelector
	commonLabels       map[string]string
	rule               model.Rule
	active             map[uint64]*Alert
	state              AlertState
}

type Fire struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

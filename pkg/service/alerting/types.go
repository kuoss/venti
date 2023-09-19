package alerting

import (
	"fmt"
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

func (s AlertState) String() string {
	switch s {
	case StateInactive:
		return "inactive"
	case StatePending:
		return "pending"
	case StateFiring:
		return "firing"
	}
	panic(fmt.Errorf("unknown alert state: %d", s))
}

type Alert struct {
	State     AlertState
	CreatedAt time.Time
	UpdatedAt time.Time

	Labels      map[string]string
	Annotations map[string]string
}

type AlertingRuleGroup struct {
	datasourceSelector model.DatasourceSelector
	groupLabels        map[string]string
	alertingRules      []AlertingRule
}

type AlertingRule struct {
	rule   model.Rule
	active map[uint64]*Alert
}

func (r AlertingRule) State() AlertState {
	maxState := StateInactive
	for _, alert := range r.active {
		if alert.State > maxState {
			maxState = alert.State
		}
	}
	return maxState
}

type Fire struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

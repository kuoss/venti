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
	State       AlertState        `json:"state"`
	CreatedAt   time.Time         `json:"createdAt,omitempty"`
	UpdatedAt   time.Time         `json:"updatedAt,omitempty"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

type AlertingRuleGroup struct {
	DatasourceSelector model.DatasourceSelector `json:"datasourceSelector,omitempty"`
	GroupLabels        map[string]string        `json:"groupLabels,omitempty"`
	AlertingRules      []AlertingRule           `json:"alertingRules,omitempty"`
}

type AlertingRule struct {
	Rule   model.Rule        `json:"rule,omitempty"`
	Active map[uint64]*Alert `json:"active,omitempty"`
}

func (r AlertingRule) State() AlertState {
	maxState := StateInactive
	for _, alert := range r.Active {
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

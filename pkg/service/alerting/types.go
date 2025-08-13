package alerting

import (
	"fmt"
	"time"

	"github.com/kuoss/venti/pkg/model"
)

// https://github.com/prometheus/prometheus/blob/v3.5.0/rules/alerting.go

// AlertState denotes the state of an active alert.
type AlertState int

const (
	// StateInactive is the state of an alert that is neither firing nor pending.
	StateInactive AlertState = iota
	// StatePending is the state of an alert that has been active for less than
	// the configured threshold duration.
	StatePending
	// StateFiring is the state of an alert that has been active for longer than
	// the configured threshold duration.
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

// modified
type Alert struct {
	State AlertState `json:"state"`

	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`

	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
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

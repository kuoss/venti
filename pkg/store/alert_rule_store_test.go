package store

import (
	"os"
	"testing"

	"github.com/kuoss/venti/pkg/model"
	"github.com/stretchr/testify/assert"
)

var (
	sampleRuleFiles []model.RuleFile
)

func init() {
	_ = os.Chdir("../..")
	sampleRuleFiles = []model.RuleFile{model.RuleFile{
		Kind:               "AlertRuleFile",
		CommonLabels:       map[string]string{"severity": "silence"},
		DatasourceSelector: model.DatasourceSelector{Type: model.DatasourceTypePrometheus},
		Groups: []model.RuleGroup{{
			Name:     "sample",
			Interval: 0,
			Limit:    0,
			Rules: []model.Rule{
				{Alert: "S00-AlwaysOn", Expr: "vector(1234)", For: 0, Annotations: map[string]string{"summary": "AlwaysOn value={{ $value }}"}},
				{Alert: "S01-Monday", Expr: "day_of_week() == 1 and hour() < 2", For: 0, Annotations: map[string]string{"summary": "Monday"}},
				{Alert: "S02-NewNamespace", Expr: "time() - kube_namespace_created < 120", For: 0, Annotations: map[string]string{"summary": "labels={{ $labels }} namespace={{ $labels.namespace }} value={{ $value }}"}},
			}},
		}}}
}

func TestNewAlertRuleStore(t *testing.T) {
	s, err := NewAlertRuleStore("")
	assert.Nil(t, err)
	assert.Equal(t, &AlertRuleStore{alertRuleFiles: sampleRuleFiles}, s)
}

func TestAlertRuleFiles(t *testing.T) {
	s, err := NewAlertRuleStore("")
	assert.Nil(t, err)
	assert.Equal(t, sampleRuleFiles, s.AlertRuleFiles())
}

package alertrule

import (
	"fmt"
	"os"
	"testing"

	"github.com/kuoss/venti/pkg/model"
	"github.com/stretchr/testify/assert"
)

var (
	alertRuleFiles = []model.RuleFile{{
		Kind:               "AlertRuleFile",
		CommonLabels:       map[string]string{"rulefile": "sample-v3", "severity": "silence"},
		DatasourceSelector: model.DatasourceSelector{System: "", Type: "prometheus"},
		RuleGroups: []model.RuleGroup{
			{Name: "sample", Interval: 0, Limit: 0, Rules: []model.Rule{
				{Record: "", Alert: "S00-AlwaysOn", Expr: "vector(1234)", For: 0, KeepFiringFor: 0, Labels: map[string]string{"hello": "world"}, Annotations: map[string]string{"summary": "AlwaysOn value={{ $value }}"}},
				{Record: "", Alert: "S01-Monday", Expr: "day_of_week() == 1 and hour() < 2", For: 0, KeepFiringFor: 0, Labels: map[string]string(nil), Annotations: map[string]string{"summary": "Monday"}},
				{Record: "", Alert: "S02-NewNamespace", Expr: "time() - kube_namespace_created < 120", For: 0, KeepFiringFor: 0, Labels: map[string]string(nil), Annotations: map[string]string{"summary": "labels={{ $labels }} namespace={{ $labels.namespace }} value={{ $value }}"}},
			}}}}}
)

func init() {
	_ = os.Chdir("../../..")

}

func TestNew(t *testing.T) {
	testCases := []struct {
		pattern   string
		want      *AlertRuleStore
		wantError string
	}{
		// ok
		{
			"etc/alertrules/*.yaml",
			&AlertRuleStore{alertRuleFiles: alertRuleFiles},
			"",
		},
		{
			"",
			&AlertRuleStore{alertRuleFiles: []model.RuleFile(nil)},
			"",
		},
		// error
		{
			"[]",
			(*AlertRuleStore)(nil),
			"error on Glob: syntax error in pattern",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(tt *testing.T) {
			store, err := New(tc.pattern)
			if tc.wantError == "" {
				assert.NoError(tt, err)
			} else {
				assert.EqualError(tt, err, tc.wantError)
			}
			assert.Equal(tt, tc.want, store)
		})
	}
}

func TestAlertRuleFiles(t *testing.T) {
	testCases := []struct {
		pattern string
		want    []model.RuleFile
	}{
		{
			"etc/alertrules/*.yaml",
			alertRuleFiles,
		},
		{
			"asdf/asdf.yaml",
			[]model.RuleFile(nil),
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(tt *testing.T) {
			store, err := New(tc.pattern)
			assert.Nil(t, err)
			ruleFiles := store.AlertRuleFiles()
			assert.Equal(tt, tc.want, ruleFiles)
		})
	}
}

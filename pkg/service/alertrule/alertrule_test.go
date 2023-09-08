package alertrule

import (
	"fmt"
	"os"
	"testing"

	"github.com/kuoss/venti/pkg/model"
	"github.com/stretchr/testify/assert"
)

var (
	ruleFiles = []model.RuleFile{{
		Kind:               "AlertRuleFile",
		CommonLabels:       map[string]string{"rulefile": "sample-v3", "severity": "silence"},
		DatasourceSelector: model.DatasourceSelector{System: "", Type: "prometheus"},
		RuleGroups: []model.RuleGroup{
			{Name: "sample", Interval: 0, Limit: 0, Rules: []model.Rule{
				{Record: "", Alert: "S00-AlwaysOn", Expr: "vector(1234)", For: 0, KeepFiringFor: 0, Labels: map[string]string{"hello": "world"}, Annotations: map[string]string{"summary": "AlwaysOn value={{ $value }}"}},
				{Record: "", Alert: "S01-Monday", Expr: "day_of_week() == 1 and hour() < 2", For: 0, KeepFiringFor: 0, Labels: map[string]string(nil), Annotations: map[string]string{"summary": "Monday"}},
				{Record: "", Alert: "S02-NewNamespace", Expr: "time() - kube_namespace_created < 120", For: 0, KeepFiringFor: 0, Labels: map[string]string(nil), Annotations: map[string]string{"summary": "labels={{ $labels }} namespace={{ $labels.namespace }} value={{ $value }}"}},
			}}}}}
	ruleFiles1 = []model.RuleFile{
		{Kind: "AlertRuleFile", CommonLabels: map[string]string{"rulefile": "sample-v3", "severity": "silence"}, DatasourceSelector: model.DatasourceSelector{System: "", Type: "prometheus"}, RuleGroups: []model.RuleGroup{
			{Name: "sample", Interval: 0, Limit: 0, Rules: []model.Rule{
				{Record: "", Alert: "S00-AlwaysOn", Expr: "vector(1234)", For: 0, KeepFiringFor: 0, Labels: map[string]string{"hello": "world"}, Annotations: map[string]string{"summary": "AlwaysOn value={{ $value }}"}},
				{Record: "", Alert: "S01-Monday", Expr: "day_of_week() == 1 and hour() < 2", For: 0, KeepFiringFor: 0, Labels: map[string]string(nil), Annotations: map[string]string{"summary": "Monday"}},
				{Record: "", Alert: "S02-NewNamespace", Expr: "time() - kube_namespace_created < 120", For: 0, KeepFiringFor: 0, Labels: map[string]string(nil), Annotations: map[string]string{"summary": "labels={{ $labels }} namespace={{ $labels.namespace }} value={{ $value }}"}}}}}}}
)

func init() {
	err := os.Chdir("../../..")
	if err != nil {
		panic(err)
	}
}

func TestNew(t *testing.T) {
	testCases := []struct {
		pattern   string
		want      *AlertRuleService
		wantError string
	}{
		// ok
		{
			"etc/alertrules/*.y*ml",
			&AlertRuleService{AlertRuleFiles: ruleFiles1},
			"",
		},
		{
			"",
			&AlertRuleService{AlertRuleFiles: ruleFiles1},
			"",
		},
		// error
		{
			"asdf",
			&AlertRuleService{AlertRuleFiles: []model.RuleFile(nil)},
			"",
		},
		{
			"[]",
			(*AlertRuleService)(nil),
			"error on Glob: syntax error in pattern",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			service, err := New(tc.pattern)
			if tc.wantError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantError)
			}
			assert.Equal(t, tc.want, service)
		})
	}
}

func TestAlertRuleFiles(t *testing.T) {
	testCases := []struct {
		pattern string
		want    []model.RuleFile
	}{
		{
			"",
			ruleFiles,
		},
		{
			"asdf",
			[]model.RuleFile(nil),
		},
		{
			"etc/alertrules/*.yml",
			ruleFiles,
		},
		{
			"etc/alertrules/*.yaml",
			[]model.RuleFile(nil),
		},
		{
			"etc/alertrules/*.y*ml",
			ruleFiles,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			service, err := New(tc.pattern)
			assert.NoError(t, err)
			ruleFiles := service.GetAlertRuleFiles()
			assert.Equal(t, tc.want, ruleFiles)
		})
	}
}

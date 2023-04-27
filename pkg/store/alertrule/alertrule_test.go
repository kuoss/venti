package alertrule

import (
	"fmt"
	"os"
	"runtime"
	"strings"
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

func line() string {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Sprintf("%s:%d", file[strings.LastIndex(file, "/")+1:], line)
}
func tcname(i int, line string) string {
	return fmt.Sprintf("TESTCASE#%d:%s", i, line)
}

func TestNew(t *testing.T) {
	testCases := []struct {
		line      string
		pattern   string
		want      *AlertRuleStore
		wantError string
	}{
		// ok
		{
			line(),
			"etc/alertrules/*.y*ml",
			&AlertRuleStore{alertRuleFiles: alertRuleFiles},
			"",
		},
		{
			line(),
			"",
			&AlertRuleStore{alertRuleFiles: alertRuleFiles},
			"",
		},
		// error
		{
			line(),
			"asdf",
			&AlertRuleStore{alertRuleFiles: []model.RuleFile(nil)},
			"",
		},
		{
			line(),
			"[]",
			(*AlertRuleStore)(nil),
			"error on Glob: syntax error in pattern",
		},
	}
	for i, tc := range testCases {
		t.Run(tcname(i, tc.line), func(t *testing.T) {
			store, err := New(tc.pattern)
			if tc.wantError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantError)
			}
			assert.Equal(t, tc.want, store)
		})
	}
}

func TestAlertRuleFiles(t *testing.T) {
	testCases := []struct {
		line    string
		pattern string
		want    []model.RuleFile
	}{
		{
			line(),
			"",
			alertRuleFiles,
		},
		{
			line(),
			"asdf",
			[]model.RuleFile(nil),
		},
		{
			line(),
			"etc/alertrules/*.yml",
			alertRuleFiles,
		},
		{
			line(),
			"etc/alertrules/*.yaml",
			[]model.RuleFile(nil),
		},
		{
			line(),
			"etc/alertrules/*.y*ml",
			alertRuleFiles,
		},
	}
	for i, tc := range testCases {
		t.Run(tcname(i, tc.line), func(t *testing.T) {
			store, err := New(tc.pattern)
			assert.Nil(t, err)
			ruleFiles := store.AlertRuleFiles()
			assert.Equal(t, tc.want, ruleFiles)
		})
	}
}

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
				{Record: "", Alert: "PodNotHealthy", Expr: "sum by (namespace, pod) (kube_pod_status_phase{phase=~\"Pending|Unknown|Failed\"}) > 0", For: 3000000000, KeepFiringFor: 0, Labels: map[string]string(nil), Annotations: map[string]string{"summary": "{{ $labels.namespace }}/{{ $labels.pod }}"}},
			}}}}}
	ruleFiles1 = []model.RuleFile{
		{Kind: "AlertRuleFile", CommonLabels: map[string]string{"rulefile": "sample-v3", "severity": "silence"}, DatasourceSelector: model.DatasourceSelector{System: "", Type: "prometheus"}, RuleGroups: []model.RuleGroup{
			{Name: "sample", Interval: 0, Limit: 0, Rules: []model.Rule{
				{Record: "", Alert: "S00-AlwaysOn", Expr: "vector(1234)", For: 0, KeepFiringFor: 0, Labels: map[string]string{"hello": "world"}, Annotations: map[string]string{"summary": "AlwaysOn value={{ $value }}"}},
				{Record: "", Alert: "S01-Monday", Expr: "day_of_week() == 1 and hour() < 2", For: 0, KeepFiringFor: 0, Labels: map[string]string(nil), Annotations: map[string]string{"summary": "Monday"}},
				{Record: "", Alert: "S02-NewNamespace", Expr: "time() - kube_namespace_created < 120", For: 0, KeepFiringFor: 0, Labels: map[string]string(nil), Annotations: map[string]string{"summary": "labels={{ $labels }} namespace={{ $labels.namespace }} value={{ $value }}"}},
				{Record: "", Alert: "PodNotHealthy", Expr: "sum by (namespace, pod) (kube_pod_status_phase{phase=~\"Pending|Unknown|Failed\"}) > 0", For: 3000000000, KeepFiringFor: 0, Labels: map[string]string(nil), Annotations: map[string]string{"summary": "{{ $labels.namespace }}/{{ $labels.pod }}"}},
			}}}}}
)

func init() {
	err := os.Chdir("../../..")
	if err != nil {
		panic(err)
	}
}

func TestNew(t *testing.T) {
	testCases := []struct {
		name      string
		pattern   string
		want      *AlertRuleService
		wantError string
	}{
		// ok
		{
			"ok",
			"etc/alertrules/*.y*ml",
			&AlertRuleService{AlertRuleFiles: ruleFiles1},
			"",
		},
		{
			"ok",
			"",
			&AlertRuleService{AlertRuleFiles: ruleFiles1},
			"",
		},
		// error
		{
			"error",
			"asdf",
			&AlertRuleService{AlertRuleFiles: []model.RuleFile(nil)},
			"",
		},
		{
			"error",
			"[]",
			(*AlertRuleService)(nil),
			"glob err: syntax error in pattern",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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

func TestNew_tempFiles(t *testing.T) {
	testCases := []struct {
		filename  string
		content   string
		want      *AlertRuleService
		wantError string
	}{
		{
			"test.ok.yaml",
			`
groups:
- name: info
  rules:
  - alert: hello
    expr: greet > 90
    for: 1m
    annotations:
      summary: "hello world"
`,
			&AlertRuleService{AlertRuleFiles: []model.RuleFile{{RuleGroups: []model.RuleGroup{
				{Name: "info", Rules: []model.Rule{{
					Alert:       "hello",
					Expr:        "greet > 90",
					For:         60000000000,
					Annotations: map[string]string{"summary": "hello world"}}}}}}}},
			"",
		},
		{
			"test.err.yaml",
			`
groups:
- name: info
  rules:
  - alert: hello
    expr: greet > 90
    for: 0m
    annotations:
      summary: "hello" world"
`,
			nil,
			"loadAlertRuleFileFromFilename err: unmarshalStrict err: yaml: line 8: did not find expected key",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			os.WriteFile(tc.filename, []byte(tc.content), 0660)
			defer func() {
				os.RemoveAll(tc.filename)
			}()
			service, err := New(tc.filename)
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

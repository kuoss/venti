package store

import (
	"os"
	"testing"

	"github.com/prometheus/prometheus/model/rulefmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func init() {
	_ = os.Chdir("../..")
}

func TestNewAlertRuleStore(t *testing.T) {
	ars, err := NewAlertRuleStore("")
	assert.Nil(t, err)
	assert.Equal(t, ars, &AlertRuleStore{ruleGroupsList: []rulefmt.RuleGroups{
		{Groups: []rulefmt.RuleGroup{
			{Name: "sample", Interval: 0, Limit: 0, Rules: []rulefmt.RuleNode{
				{Record: yaml.Node{Kind: 0x0, Style: 0x0, Tag: "", Value: "", Anchor: "", Alias: (*yaml.Node)(nil), Content: []*yaml.Node(nil), HeadComment: "", LineComment: "", FootComment: "", Line: 0, Column: 0}, Alert: yaml.Node{Kind: 0x8, Style: 0x0, Tag: "!!str", Value: "S00-AlwaysOn", Anchor: "", Alias: (*yaml.Node)(nil), Content: []*yaml.Node(nil), HeadComment: "", LineComment: "", FootComment: "", Line: 7, Column: 12}, Expr: yaml.Node{Kind: 0x8, Style: 0x0, Tag: "!!str", Value: "vector(1234)", Anchor: "", Alias: (*yaml.Node)(nil), Content: []*yaml.Node(nil), HeadComment: "", LineComment: "", FootComment: "", Line: 8, Column: 11}, For: 0, Labels: map[string]string(nil), Annotations: map[string]string{"summary": "AlwaysOn value={{ $value }}"}},
				{Record: yaml.Node{Kind: 0x0, Style: 0x0, Tag: "", Value: "", Anchor: "", Alias: (*yaml.Node)(nil), Content: []*yaml.Node(nil), HeadComment: "", LineComment: "", FootComment: "", Line: 0, Column: 0}, Alert: yaml.Node{Kind: 0x8, Style: 0x0, Tag: "!!str", Value: "S01-Monday", Anchor: "", Alias: (*yaml.Node)(nil), Content: []*yaml.Node(nil), HeadComment: "", LineComment: "", FootComment: "", Line: 12, Column: 12}, Expr: yaml.Node{Kind: 0x8, Style: 0x0, Tag: "!!str", Value: "day_of_week() == 1 and hour() < 2", Anchor: "", Alias: (*yaml.Node)(nil), Content: []*yaml.Node(nil), HeadComment: "", LineComment: "", FootComment: "", Line: 13, Column: 11}, For: 0, Labels: map[string]string(nil), Annotations: map[string]string{"summary": "Monday"}},
				{Record: yaml.Node{Kind: 0x0, Style: 0x0, Tag: "", Value: "", Anchor: "", Alias: (*yaml.Node)(nil), Content: []*yaml.Node(nil), HeadComment: "", LineComment: "", FootComment: "", Line: 0, Column: 0}, Alert: yaml.Node{Kind: 0x8, Style: 0x0, Tag: "!!str", Value: "S02-NewNamespace", Anchor: "", Alias: (*yaml.Node)(nil), Content: []*yaml.Node(nil), HeadComment: "", LineComment: "", FootComment: "", Line: 17, Column: 12}, Expr: yaml.Node{Kind: 0x8, Style: 0x0, Tag: "!!str", Value: "time() - kube_namespace_created < 120", Anchor: "", Alias: (*yaml.Node)(nil), Content: []*yaml.Node(nil), HeadComment: "", LineComment: "", FootComment: "", Line: 18, Column: 11}, For: 0, Labels: map[string]string(nil), Annotations: map[string]string{"summary": "labels={{ $labels }} namespace={{ $labels.namespace }} value={{ $value }}"}},
			}}}}}})
}

func TestRuleGroupsList(t *testing.T) {
	ars, err := NewAlertRuleStore("")
	assert.Nil(t, err)
	assert.Equal(t, ars.RuleGroupsList(), []rulefmt.RuleGroups{
		{Groups: []rulefmt.RuleGroup{
			{Name: "sample", Interval: 0, Limit: 0, Rules: []rulefmt.RuleNode{
				{Record: yaml.Node{Kind: 0x0, Style: 0x0, Tag: "", Value: "", Anchor: "", Alias: (*yaml.Node)(nil), Content: []*yaml.Node(nil), HeadComment: "", LineComment: "", FootComment: "", Line: 0, Column: 0}, Alert: yaml.Node{Kind: 0x8, Style: 0x0, Tag: "!!str", Value: "S00-AlwaysOn", Anchor: "", Alias: (*yaml.Node)(nil), Content: []*yaml.Node(nil), HeadComment: "", LineComment: "", FootComment: "", Line: 7, Column: 12}, Expr: yaml.Node{Kind: 0x8, Style: 0x0, Tag: "!!str", Value: "vector(1234)", Anchor: "", Alias: (*yaml.Node)(nil), Content: []*yaml.Node(nil), HeadComment: "", LineComment: "", FootComment: "", Line: 8, Column: 11}, For: 0, Labels: map[string]string(nil), Annotations: map[string]string{"summary": "AlwaysOn value={{ $value }}"}},
				{Record: yaml.Node{Kind: 0x0, Style: 0x0, Tag: "", Value: "", Anchor: "", Alias: (*yaml.Node)(nil), Content: []*yaml.Node(nil), HeadComment: "", LineComment: "", FootComment: "", Line: 0, Column: 0}, Alert: yaml.Node{Kind: 0x8, Style: 0x0, Tag: "!!str", Value: "S01-Monday", Anchor: "", Alias: (*yaml.Node)(nil), Content: []*yaml.Node(nil), HeadComment: "", LineComment: "", FootComment: "", Line: 12, Column: 12}, Expr: yaml.Node{Kind: 0x8, Style: 0x0, Tag: "!!str", Value: "day_of_week() == 1 and hour() < 2", Anchor: "", Alias: (*yaml.Node)(nil), Content: []*yaml.Node(nil), HeadComment: "", LineComment: "", FootComment: "", Line: 13, Column: 11}, For: 0, Labels: map[string]string(nil), Annotations: map[string]string{"summary": "Monday"}},
				{Record: yaml.Node{Kind: 0x0, Style: 0x0, Tag: "", Value: "", Anchor: "", Alias: (*yaml.Node)(nil), Content: []*yaml.Node(nil), HeadComment: "", LineComment: "", FootComment: "", Line: 0, Column: 0}, Alert: yaml.Node{Kind: 0x8, Style: 0x0, Tag: "!!str", Value: "S02-NewNamespace", Anchor: "", Alias: (*yaml.Node)(nil), Content: []*yaml.Node(nil), HeadComment: "", LineComment: "", FootComment: "", Line: 17, Column: 12}, Expr: yaml.Node{Kind: 0x8, Style: 0x0, Tag: "!!str", Value: "time() - kube_namespace_created < 120", Anchor: "", Alias: (*yaml.Node)(nil), Content: []*yaml.Node(nil), HeadComment: "", LineComment: "", FootComment: "", Line: 18, Column: 11}, For: 0, Labels: map[string]string(nil), Annotations: map[string]string{"summary": "labels={{ $labels }} namespace={{ $labels.namespace }} value={{ $value }}"}},
			}}}}})
}

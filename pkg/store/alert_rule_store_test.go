package store

import (
	"os"
	"testing"

	"github.com/kuoss/venti/pkg/model"
	"github.com/stretchr/testify/assert"
)

func init() {
	_ = os.Chdir("../..")
}

func TestNewAlertRuleStore(t *testing.T) {
	ars, err := NewAlertRuleStore("")
	assert.Nil(t, err)
	assert.Equal(t, &AlertRuleStore{
		ruleGroupsList: []model.RuleGroups{{
			Groups: []model.RuleGroup{{Name: "sample", Interval: 0, Limit: 0, Rules: []model.Rule{}}},
		}},
	}, ars)
}

func TestRuleGroupsList(t *testing.T) {
	ars, err := NewAlertRuleStore("")
	assert.Nil(t, err)
	assert.Equal(t, []model.RuleGroups{{Groups: []model.RuleGroup{{Name: "sample", Interval: 0, Limit: 0, Rules: []model.Rule{}}}}}, ars.RuleGroupsList())
}

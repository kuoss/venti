package store

import (
	"os"
	"testing"

	"github.com/kuoss/venti/pkg/mock"
	"github.com/kuoss/venti/pkg/model"
	"github.com/stretchr/testify/assert"
)

var (
	alertRuleFiles []model.RuleFile
)

func init() {
	_ = os.Chdir("../..")
	alertRuleFiles = mock.AlertRuleFiles()
}

func TestNewAlertRuleStore(t *testing.T) {
	s, err := NewAlertRuleStore("etc/alertrules/*.yaml")
	assert.Nil(t, err)
	assert.Equal(t, &AlertRuleStore{alertRuleFiles: alertRuleFiles}, s)
}

func TestAlertRuleFiles(t *testing.T) {
	s, err := NewAlertRuleStore("etc/alertrules/*.yaml")
	assert.Nil(t, err)
	assert.Equal(t, alertRuleFiles, s.AlertRuleFiles())
}

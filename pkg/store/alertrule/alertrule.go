package alertrule

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/venti/pkg/model"
	"gopkg.in/yaml.v2"
)

type AlertRuleStore struct {
	alertRuleFiles []model.RuleFile
}

func New(pattern string) (alertRuleStore *AlertRuleStore, err error) {
	logger.Infof("loading alertrules...")
	if pattern == "" {
		pattern = "etc/alertrules/*.y*ml"
	}
	files, err := filepath.Glob(pattern)
	if err != nil {
		err = fmt.Errorf("error on Glob: %w", err)
		return
	}
	var alertRuleFiles []model.RuleFile
	for _, filename := range files {
		alertRuleFile, err := loadAlertRuleFileFromFilename(filename)
		if err != nil {
			logger.Warnf("error on loadAlertRuleGroupsFromFile(skipped): %s", err)
			continue
		}
		alertRuleFiles = append(alertRuleFiles, *alertRuleFile)
	}
	alertRuleStore = &AlertRuleStore{alertRuleFiles: alertRuleFiles}
	return
}

func loadAlertRuleFileFromFilename(filename string) (*model.RuleFile, error) {
	logger.Infof("load alertrule file: %s", filename)
	yamlBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error on ReadFile: %w", err)
	}
	var alertRuleFile *model.RuleFile
	if err := yaml.UnmarshalStrict(yamlBytes, &alertRuleFile); err != nil {
		return nil, fmt.Errorf("error on UnmarshalStrict: %w", err)
	}
	return alertRuleFile, nil
}

func (s *AlertRuleStore) AlertRuleFiles() []model.RuleFile {
	return s.alertRuleFiles
}

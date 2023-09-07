package alertrule

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/venti/pkg/model"
	"gopkg.in/yaml.v2"
)

type AlertRuleService struct {
	AlertRuleFiles []model.RuleFile
}

func New(pattern string) (alertRuleService *AlertRuleService, err error) {
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
			logger.Warnf("loadAlertRuleFileFromFilename err: %s", err)
			continue
		}
		alertRuleFiles = append(alertRuleFiles, *alertRuleFile)
	}
	alertRuleService = &AlertRuleService{AlertRuleFiles: alertRuleFiles}
	return
}

func loadAlertRuleFileFromFilename(filename string) (*model.RuleFile, error) {
	logger.Infof("load alertrule file: %s", filename)
	yamlBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("readFile err: %w", err)
	}
	var alertRuleFile *model.RuleFile
	if err := yaml.UnmarshalStrict(yamlBytes, &alertRuleFile); err != nil {
		return nil, fmt.Errorf("unmarshalStrict err: %w", err)
	}
	return alertRuleFile, nil
}

func (s *AlertRuleService) GetAlertRuleFiles() []model.RuleFile {
	return s.AlertRuleFiles
}

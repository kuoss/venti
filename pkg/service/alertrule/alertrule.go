package alertrule

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/util"
)

type AlertRuleService struct {
	AlertRuleFiles []model.RuleFile
}

func New(pattern string) (*AlertRuleService, error) {
	logger.Infof("loading alertrules...")
	if pattern == "" {
		pattern = "etc/alertrules/*.y*ml"
	}
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("glob err: %w", err)
	}
	var alertRuleFiles []model.RuleFile
	for _, filename := range files {
		alertRuleFile, err := loadAlertRuleFileFromFilename(filename)
		if err != nil {
			return nil, fmt.Errorf("loadAlertRuleFileFromFilename err: %w", err)
		}
		alertRuleFiles = append(alertRuleFiles, *alertRuleFile)
	}
	return &AlertRuleService{AlertRuleFiles: alertRuleFiles}, nil
}

func loadAlertRuleFileFromFilename(filename string) (*model.RuleFile, error) {
	logger.Infof("load alertrule file: %s", filename)
	yamlBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("readFile err: %w", err)
	}
	var alertRuleFile *model.RuleFile
	if err := util.UnmarshalStrict(yamlBytes, &alertRuleFile); err != nil {
		return nil, fmt.Errorf("unmarshalStrict err: %w", err)
	}
	return alertRuleFile, nil
}

func (s *AlertRuleService) GetAlertRuleFiles() []model.RuleFile {
	return s.AlertRuleFiles
}

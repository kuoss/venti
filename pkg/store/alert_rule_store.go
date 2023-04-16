package store

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kuoss/venti/pkg/model"
	"gopkg.in/yaml.v2"
)

type AlertRuleStore struct {
	alertRuleFiles []model.RuleFile
}

func NewAlertRuleStore(pattern string) (*AlertRuleStore, error) {
	log.Println("Loading alertRules...")
	files, err := filepath.Glob("etc/alertrules/*.yaml")
	if err != nil {
		return nil, err
	}
	var alertRuleFiles []model.RuleFile
	for _, filename := range files {
		alertRuleFile, err := loadAlertRuleFileFromFilename(filename)
		if err != nil {
			log.Printf("Warning: error on loadAlertRuleGroupsFromFile(skipped): %s", err)
			continue
		}
		alertRuleFiles = append(alertRuleFiles, *alertRuleFile)
	}
	return &AlertRuleStore{alertRuleFiles: alertRuleFiles}, nil
}

func loadAlertRuleFileFromFilename(filename string) (*model.RuleFile, error) {
	log.Printf("load alertrule file: %s\n", filename)
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

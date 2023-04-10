package store

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/kuoss/venti/pkg/model"
	"gopkg.in/yaml.v2"
)

type AlertRuleStore struct {
	ruleGroupsList []model.RuleGroups
}

func NewAlertRuleStore(pattern string) (*AlertRuleStore, error) {
	log.Println("Loading alertRules...")
	files, err := filepath.Glob("etc/alertrules/*.yaml")
	if err != nil {
		return nil, err
	}
	var ruleGroupsList []model.RuleGroups
	for _, filename := range files {
		ruleGroups, err := loadAlertRuleGroupsFromFile(filename)
		if err != nil {
			log.Printf("Warning: error on loadAlertRuleGroupsFromFile(skipped): %s", err)
			continue
		}
		ruleGroupsList = append(ruleGroupsList, *ruleGroups)
	}
	return &AlertRuleStore{ruleGroupsList: ruleGroupsList}, nil
}

func loadAlertRuleGroupsFromFile(filename string) (*model.RuleGroups, error) {
	log.Printf("load alertrules file: %s\n", filename)
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error on Open: %w", err)
	}
	var ruleGroups *model.RuleGroups
	yamlBytes, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("error on ReadAll: %w", err)
	}
	if err := yaml.UnmarshalStrict(yamlBytes, &ruleGroups); err != nil {
		return nil, fmt.Errorf("error on UnmarshalStrict: %w", err)
	}
	return ruleGroups, nil
}

func (s *AlertRuleStore) RuleGroupsList() []model.RuleGroups {
	return s.ruleGroupsList
}

package store

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/prometheus/prometheus/model/rulefmt"
)

type AlertRuleStore struct {
	ruleGroupsList []rulefmt.RuleGroups
}

func NewAlertRuleStore(pattern string) (*AlertRuleStore, error) {
	log.Println("Loading alertRules...")
	files, err := filepath.Glob("etc/alertrules/*.yaml")
	if err != nil {
		return nil, err
	}
	var ruleGroupsList []rulefmt.RuleGroups
	for _, filename := range files {
		log.Printf("alertRule file: %s\n", filename)
		f, err := os.Open(filename)
		if err != nil {
			return nil, fmt.Errorf("error on Open: %w", err)
		}
		var ruleGroups *rulefmt.RuleGroups
		err = loadYaml(f, &ruleGroups)
		if err != nil {
			return nil, fmt.Errorf("error on loadYaml: %w", err)
		}
		ruleGroupsList = append(ruleGroupsList, *ruleGroups)
	}
	return &AlertRuleStore{ruleGroupsList: ruleGroupsList}, nil
}

func (ars *AlertRuleStore) RuleGroupsList() []rulefmt.RuleGroups {
	return ars.ruleGroupsList
}

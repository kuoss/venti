package store

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/prometheus/prometheus/model/rulefmt"
)

type AlertRuleStore struct {
	ruleGroupsSlice []rulefmt.RuleGroups
}

func NewAlertRuleStore(pattern string) (*AlertRuleStore, error) {
	log.Println("Loading alertRules...")
	files, err := filepath.Glob("etc/alertrules/*.yaml")
	if err != nil {
		return nil, err
	}
	var ruleGroupsSlice []rulefmt.RuleGroups
	for _, filename := range files {
		log.Printf("alertRuleGroup file: %s\n", filename)
		f, err := os.Open(filename)
		if err != nil {
			return nil, fmt.Errorf("error on Open: %w", err)
		}
		var ruleGroups *rulefmt.RuleGroups
		err = loadYaml(f, &ruleGroups)
		if err != nil {
			return nil, fmt.Errorf("error on loadYaml: %w", err)
		}
		ruleGroupsSlice = append(ruleGroupsSlice, *ruleGroups)
	}
	return &AlertRuleStore{ruleGroupsSlice: ruleGroupsSlice}, nil
}

func (ars *AlertRuleStore) GroupsSlice() []rulefmt.RuleGroups {
	return ars.ruleGroupsSlice
}

package store

import (
	"log"
	"os"
	"path/filepath"

	"github.com/prometheus/prometheus/model/rulefmt"
)

type AlertRuleStore struct {
	ruleGroups rulefmt.RuleGroups
}

func NewAlertRuleStore(pattern string) (*AlertRuleStore, error) {
	log.Println("Loading alertrules...")
	// default: "etc/alertrules/*.yaml"
	files, err := filepath.Glob("etc/alertrules/*.yaml")
	if err != nil {
		return nil, err
	}

	var alertRuleGroups rulefmt.RuleGroups
	for _, filename := range files {
		log.Printf("alertrulegroup file: %s\n", filename)
		f, err := os.Open(filename)
		if err != nil {
			log.Printf("error open alertrulegroup file: %s\n", err.Error())
		}
		var rg *rulefmt.RuleGroup = &rulefmt.RuleGroup{}
		err = loadYaml(f, rg)
		if err != nil {
			return nil, err
		}
		alertRuleGroups.Groups = append(alertRuleGroups.Groups, *rg)
	}
	return &AlertRuleStore{ruleGroups: alertRuleGroups}, nil
}

func (ars *AlertRuleStore) Groups() rulefmt.RuleGroups {
	return ars.ruleGroups
}

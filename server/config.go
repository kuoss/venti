package server

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	config          Config
	datasourceStore DatasourceStore
)

func LoadConfig(version string) {
	var err error
	log.Println("Loading configurations...")
	config.Version = version
	err = loadUsersConfig()
	if err != nil {
		log.Printf("error on loadUsersConfig: %s", err)
	}
	err = loadDatasourcesConfig()
	if err != nil {
		log.Printf("error on loadDatasourcesConfig: %s", err)
	}
	datasourceStore = DatasourceStore.New(config.DatasourcesConfig)
	loadDashboards()
	loadAlertRuleGroups()
}

func GetConfig() Config {
	return config
}

func loadUsersConfig() error {
	yamlBytes, err := os.ReadFile("etc/users.yaml")
	if err != nil {
		return fmt.Errorf("cannot ReadFile: %w", err)
	}
	if err := yaml.Unmarshal(yamlBytes, &config.EtcUsersConfig); err != nil {
		return fmt.Errorf("cannot Unmarshal: %w", err)
	}
	log.Println("Users config file loaded.")
	return nil
}

func loadDatasourcesConfig() error {
	yamlBytes, err := os.ReadFile("etc/datasources.yaml")
	if err != nil {
		return fmt.Errorf("cannot ReadFile: %w", err)
	}
	if err := yaml.Unmarshal(yamlBytes, &config.DatasourcesConfig); err != nil {
		return fmt.Errorf("cannot Unmarshal: %w", err)
	}
	log.Println("Datasources config file loaded.")
	if config.DatasourcesConfig.Discovery.AnnotationKey == "" {
		config.DatasourcesConfig.Discovery.AnnotationKey = "kuoss.org/datasource"
	}
	log.Println(config.DatasourcesConfig)
	return nil
}

func glob(root string, fn func(string) bool) []string {
	var matches []string
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if fn(path) {
			matches = append(matches, path)
		}
		return nil
	})
	return matches
}

func loadDashboards() {
	log.Println("Loading dashboards...")
	// filepaths, err := filepath.Glob("etc/dashboards/**/*.yaml")
	filepaths := glob("etc/dashboards", func(path string) bool {
		return !strings.Contains(path, "/..") && filepath.Ext(path) == ".yaml"
	})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	log.Printf("filepaths===%#v", filepaths)

	var dashboard Dashboard
	for _, filepath := range filepaths {
		yamlBytes, err := os.ReadFile(filepath)
		if err != nil {
			log.Fatal(err)
		}
		if err := yaml.Unmarshal(yamlBytes, &dashboard); err != nil {
			log.Fatal(err)
		}
		config.Dashboards = append(config.Dashboards, dashboard)
		log.Println("Dashboard config file '" + filepath + "' loaded.")
	}
}

func loadAlertRuleGroups() {
	filepaths, err := filepath.Glob("etc/alertrules/*.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var alertRuleGroups []AlertRuleGroup
	for _, filepath := range filepaths {
		yamlBytes, err := os.ReadFile(filepath)
		if err != nil {
			log.Fatal(err)
		}
		var alertRuleGroupList AlertRuleGroupList
		if err := yaml.Unmarshal(yamlBytes, &alertRuleGroupList); err != nil {
			log.Fatal(err)
		}

		alertRuleGroups = append(alertRuleGroups, alertRuleGroupList.Groups...)
		log.Println("Alert rule file '" + filepath + "' loaded.")
	}

	// attach common labels to rules
	for i, group := range alertRuleGroups {
		commonLabels := group.CommonLabels
		for j, rule := range group.Rules {
			for key, value := range commonLabels {
				if rule.Labels == nil {
					alertRuleGroups[i].Rules[j].Labels = map[string]string{key: value}
					continue
				}
				if _, exists := rule.Labels[key]; !exists {
					alertRuleGroups[i].Rules[j].Labels[key] = value
				}
			}
		}
	}
	config.AlertRuleGroups = alertRuleGroups
}

func GetAlertRuleGroups() []AlertRuleGroup {
	return config.AlertRuleGroups
}

func GetDatasources() []Datasource {
	return datasourceStore.GetDatasources()
}

func GetDefaultDatasource(typ DatasourceType) (Datasource, error) {
	return datasourceStore.GetDefaultDatasource(typ)
}

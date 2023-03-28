package configuration

import (
	"fmt"
	"github.com/kuoss/venti/server"
	"github.com/kuoss/venti/server/alert"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type config struct {
	Version           string
	EtcUsersConfig    EtcUsersConfig
	DatasourcesConfig DatasourcesConfig
	Dashboards        []Dashboard
	AlertRuleGroups   []AlertRuleGroup
}

type EtcUsersConfig struct {
	EtcUsers []EtcUser `yaml:"users"`
}

type EtcUser struct {
	Username string `yaml:"username"`
	Hash     string `yaml:"hash"`
	IsAdmin  bool   `yaml:"isAdmin,omitempty"`
}

type Dashboard struct {
	Title string `json:"title"`
	Rows  []Row  `json:"rows"`
}

type DatasourcesConfig struct {
	QueryTimeout time.Duration `json:"queryTimeout,omitempty" yaml:"queryTimeout,omitempty"`
	Datasources  []Datasource  `json:"datasources" yaml:"datasources,omitempty"`
	Discovery    Discovery     `json:"discovery,omitempty" yaml:"discovery,omitempty"`
}

func Load(version string) (*config, error) {

	log.Println("Loading configurations...")

	userConfigFile, err := os.Open("etc/users.yaml")
	if err != nil {
		return nil, err
	}
	defer userConfigFile.Close()

	var userConf *EtcUsersConfig
	err = loadConfig(userConfigFile, userConf)
	if err != nil {
		return nil, fmt.Errorf("error on loading User Config: %w", err)
	}

	dsConfigFile, err := os.Open("etc/datasources.yaml")
	if err != nil {
		return nil, err
	}
	defer dsConfigFile.Close()

	var dataSourceConfig *DatasourcesConfig
	err = loadConfig(dsConfigFile, &dataSourceConfig)
	if err != nil {
		return nil, fmt.Errorf("error on loading Datasources Config: %w", err)
	}

	dashboardfilepaths := glob("etc/dashboards", func(path string) bool {
		return !strings.Contains(path, "/..") && filepath.Ext(path) == ".yaml"
	})

	for _, path := range dashboardfilepaths {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		var dashBoard *Dashboard
		err = loadConfig(f, dashBoard)
		if err != nil {
			return nil, err
		}

	}

	loadDashboards()
	loadAlertRuleGroups()

	datasourceStore = server.NewDatasourceStore(config.DatasourcesConfig)

	return &config{
		Version:           "",
		EtcUsersConfig:    nil,
		DatasourcesConfig: nil,
		Dashboards:        nil,
		AlertRuleGroups:   nil,
	}, nil
}

func GetConfig() server.Config {
	return config
}

func loadConfig(r io.Reader, c interface{}) error {
	yamlBytes, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("cannot ReadFile: %w", err)
	}
	if err := yaml.Unmarshal(yamlBytes, c); err != nil {
		return fmt.Errorf("cannot Unmarshal: %w", err)
	}
	log.Printf("Users config file loaded.\n")
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

func loadDashboard(root string) {
	log.Println("Loading dashboards...")

	log.Printf("filepaths: %#v", filepaths)

	var dashboard server.Dashboard

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

func loadAlertRuleGroups() {
	filepaths, err := filepath.Glob("etc/alertrules/*.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var alertRuleGroups []alert.AlertRuleGroup
	for _, filepath := range filepaths {
		yamlBytes, err := os.ReadFile(filepath)
		if err != nil {
			log.Fatal(err)
		}
		var alertRuleGroupList alert.AlertRuleGroupList
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

func GetAlertRuleGroups() []alert.AlertRuleGroup {
	return config.AlertRuleGroups
}

func GetDatasources() []server.Datasource {
	return datasourceStore.GetDatasources()
}

func GetDefaultDatasource(typ server.DatasourceType) (server.Datasource, error) {
	return datasourceStore.GetDefaultDatasource(typ)
}

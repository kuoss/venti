package configuration

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/kuoss/venti/pkg/model"
	"gopkg.in/yaml.v3"
)

// Load EtcUser,DatasourceConfig files only.
// TODO each Config filepath could be parameter.
func Load(version string) (*model.Config, error) {

	log.Println("Loading configurations...")

	userConfigFile, err := os.Open("etc/users.yaml")
	if err != nil {
		return nil, err
	}
	defer userConfigFile.Close()

	var userConf model.UsersConfig
	err = loadConfig(userConfigFile, &userConf)
	if err != nil {
		return nil, fmt.Errorf("error on loading User Config: %w", err)
	}

	dsConfigFile, err := os.Open("etc/datasources.yaml")
	if err != nil {
		return nil, err
	}
	defer dsConfigFile.Close()

	var dataSourceConfig *model.DatasourcesConfig
	err = loadConfig(dsConfigFile, &dataSourceConfig)
	if err != nil {
		return nil, fmt.Errorf("error on loading Datasources Config: %w", err)
	}

	return &model.Config{
		Version:           version,
		UserConfig:        userConf,
		DatasourcesConfig: dataSourceConfig,
	}, nil
}

func loadConfig(r io.Reader, c interface{}) error {
	yamlBytes, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("cannot ReadFile: %w", err)
	}
	if err := yaml.Unmarshal(yamlBytes, c); err != nil {
		return fmt.Errorf("cannot Unmarshal: %w", err)
	}
	return nil
}

// todo annotation default value check
/*
func loadDatasourcesConfig() error {
	yamlBytes, err := os.ReadFile("etc/datasources.yaml")
	if err != nil {
		return fmt.Errorf("cannot ReadFile: %w", err)
	}
	if err := yaml.Unmarshal(yamlBytes, &Config.DatasourcesConfig); err != nil {
		return fmt.Errorf("cannot Unmarshal: %w", err)
	}
	log.Println("Datasources Config file loaded.")
	if DatasourcesConfig.Discovery.AnnotationKey == "" {
		Config.DatasourcesConfig.Discovery.AnnotationKey = "kuoss.org/datasource"
	}
	log.Println(Config.DatasourcesConfig)
	return nil
}
*/

/*
func loadDashboard(root string) {
	log.Println("Loading dashboards...")

	log.Printf("filepaths: %#v", filepaths)

	var dashboard pkg.Dashboard

	for _, filepath := range filepaths {
		yamlBytes, err := os.ReadFile(filepath)
		if err != nil {
			log.Fatal(err)
		}
		if err := yaml.Unmarshal(yamlBytes, &dashboard); err != nil {
			log.Fatal(err)
		}
		Config.Dashboards = append(Config.Dashboards, dashboard)
		log.Println("Dashboard Config file '" + filepath + "' loaded.")
	}
}


*/

/*
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

*/

/*
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
	Config.AlertRuleGroups = alertRuleGroups
}

func GetAlertRuleGroups() []alert.AlertRuleGroup {
	return Config.AlertRuleGroups
}

func GetDatasources() []pkg.Datasource {
	return datasourceStore.GetDatasources()
}

func GetDefaultDatasource(typ pkg.DatasourceType) (pkg.Datasource, error) {
	return datasourceStore.GetDefaultDatasource(typ)
}


*/

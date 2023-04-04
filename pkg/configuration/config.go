package configuration

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version           string
	UserConfig        *UsersConfig
	DatasourcesConfig *DatasourcesConfig
	//Dashboards        []Dashboard
	//AlertRuleGroups   []AlertRuleGroup
}

type UsersConfig struct {
	EtcUsers []EtcUser `yaml:"users"`
}

type EtcUser struct {
	Username string `yaml:"username"`
	Hash     string `yaml:"hash"`
	IsAdmin  bool   `yaml:"isAdmin,omitempty"`
}

const (
	DatasourceTypeNone       DatasourceType = ""
	DatasourceTypePrometheus DatasourceType = "prometheus"
	DatasourceTypeLethe      DatasourceType = "lethe"
)

type DatasourcesConfig struct {
	QueryTimeout time.Duration `json:"queryTimeout,omitempty" yaml:"queryTimeout,omitempty"`
	Datasources  []*Datasource `json:"datasources" yaml:"datasources,omitempty"`
	Discovery    Discovery     `json:"discovery,omitempty" yaml:"discovery,omitempty"`
}

type DatasourceType string

type Datasource struct {
	Type              DatasourceType `json:"type" yaml:"type"`
	Name              string         `json:"name" yaml:"name"`
	URL               string         `json:"url" yaml:"url"`
	BasicAuth         bool           `json:"basicAuth" yaml:"basicAuth"`
	BasicAuthUser     string         `json:"basicAuthUser" yaml:"basicAuthUser"`
	BasicAuthPassword string         `json:"basicAuthPassword" yaml:"basicAuthPassword"`
	IsDefault         bool           `json:"isDefault,omitempty" yaml:"isDefault,omitempty"`
	IsDiscovered      bool           `json:"isDiscovered,omitempty" yaml:"isDiscovered,omitempty"`
}

type Discovery struct {
	Enabled          bool   `json:"enabled,omitempty" yaml:"enabled,omitempty"`                   // default: false
	DefaultNamespace string `json:"defaultNamespace,omitempty" yaml:"defaultNamespace,omitempty"` // default: ''
	AnnotationKey    string `json:"annotationKey,omitempty" yaml:"annotationKey,omitempty"`       // default: kuoss.org/datasource-type
	ByNamePrometheus bool   `json:"byNamePrometheus,omitempty" yaml:"byNamePrometheus,omitempty"` // deprecated
	ByNameLethe      bool   `json:"byNameLethe,omitempty" yaml:"byNameLethe,omitempty"`           // deprecated
}

// Load EtcUser,DatasourceConfig files only.
// TODO each Config filepath could be parameter.
func Load(version string) (*Config, error) {

	log.Println("Loading configurations...")

	userConfigFile, err := os.Open("etc/users.yaml")
	if err != nil {
		return nil, err
	}
	defer userConfigFile.Close()

	var userConf *UsersConfig
	err = loadConfig(userConfigFile, &userConf)
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

	/*
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
	*/

	//loadDashboards()
	//loadAlertRuleGroups()

	//datasourceStore = pkg.NewDatasourceStore(Config.DatasourcesConfig)

	return &Config{
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

package server

import (
	"io/fs"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/client-go/kubernetes"
	// "k8s.io/client-go/rest"
)

var (
	config Config
)

func LoadConfig(version string) {
	log.Println("Loading configurations...")
	config.Version = version
	loadUsersConfig()
	loadDatasourcesConfig()
	loadDashboards()
	loadAlertRuleGroups()
}

func GetConfig() Config {
	return config
}

func loadUsersConfig() {
	yamlBytes, err := os.ReadFile("etc/users.yaml")
	if err != nil {
		log.Fatal(err)
	}
	if err := yaml.Unmarshal(yamlBytes, &config.EtcUsersConfig); err != nil {
		log.Fatal(err)
	}
	log.Println("Users config file loaded.")
}

func loadDatasourcesConfig() {
	yamlBytes, err := os.ReadFile("etc/datasources.yaml")
	if err != nil {
		log.Fatal(err)
	}
	if err := yaml.Unmarshal(yamlBytes, &config.DatasourcesConfig); err != nil {
		log.Fatal(err)
	}
	log.Println(config.DatasourcesConfig)
	log.Println("Datasources config file loaded.")
}

func glob(root string, fn func(string) bool) []string {
	var matches []string
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
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

func GetDatasources() ([]Datasource, error) {
	// config, err := rest.InClusterConfig()
	// if err != nil {
	// 	return nil, err
	// }
	// clientset, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	return nil, err
	// }

	var datasources = GetConfig().DatasourcesConfig.Datasources

	// add prometheus services
	// services, err := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	// if err != nil {
	// 	return nil, err
	// }
	// for _, service := range services.Items {
	// 	if service.Namespace == "kube-system" || service.Name != "prometheus" {
	// 		continue
	// 	}
	// 	datasources = append(datasources, Datasource{
	// 		Type:         DatasourceTypePrometheus,
	// 		Host:         service.Name + "." + service.Namespace,
	// 		Port:         9090,
	// 		IsDiscovered: true,
	// 	})
	// }
	return datasources, nil
}

func GetDefaultDatasource(dstype DatasourceType) (Datasource, error) {
	var datasource Datasource
	datasources, err := GetDatasources()
	if err != nil {
		return datasource, err
	}
	for _, ds := range datasources {
		if ds.Type == dstype {
			return ds, nil
		}
	}
	return datasource, fmt.Errorf("datasource of type %s not found", dstype)
}
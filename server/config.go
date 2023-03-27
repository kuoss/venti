package server

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	config            Config
	loadedDatasources []Datasource
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

func LoadDatasources() {
	var dc = GetConfig().DatasourcesConfig
	var datasources = dc.Datasources
	if dc.Discovery.Enabled {
		discoverdDatasources, err := discoverDatasources()
		if err != nil {
			fmt.Printf("cannot discoverDatasources: %s", err)
		} else {
			datasources = append(datasources, discoverdDatasources...)
		}
	}
	datasources = setDefaultDatasources(datasources)
	loadedDatasources = datasources
	log.Printf("loadedDatasources=%v", loadedDatasources)
}

func GetDatasources() []Datasource {
	return loadedDatasources
}

func GetDefaultDatasource(typ DatasourceType) (Datasource, error) {
	for _, ds := range loadedDatasources {
		if ds.Type == typ && ds.IsDefault {
			return ds, nil
		}
	}
	return Datasource{}, fmt.Errorf("datasource of type %s not found", typ)
}

func setDefaultDatasources(datasources []Datasource) []Datasource {
	// check & ensure to set a default datasource for each type
	existsDefaultPrometheus := false
	existsDefaultLethe := false
	for _, ds := range datasources {
		if !ds.IsDefault {
			continue
		}
		switch ds.Type {
		case DatasourceTypePrometheus:
			existsDefaultPrometheus = true
			continue
		case DatasourceTypeLethe:
			existsDefaultLethe = true
			continue
		}
	}
	if !existsDefaultPrometheus {
		for i, ds := range datasources {
			if ds.Type == DatasourceTypePrometheus {
				datasources[i].IsDefault = true
			}
		}
	}
	if !existsDefaultLethe {
		for i, ds := range datasources {
			if ds.Type == DatasourceTypeLethe {
				datasources[i].IsDefault = true
			}
		}
	}
	return datasources
}

func listServices() ([]v1.Service, error) {
	var config, err = rest.InClusterConfig()
	if err != nil {
		return []v1.Service{}, fmt.Errorf("cannot InClusterConfig: %w", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return []v1.Service{}, fmt.Errorf("cannot NewForConfig: %w", err)
	}

	services, err := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return []v1.Service{}, fmt.Errorf("cannot ListServices: %w", err)
	}
	return services.Items, nil
}

func discoverDatasources() ([]Datasource, error) {
	var datasources = []Datasource{}

	services, err := listServices()
	if err != nil {
		return datasources, fmt.Errorf("cannot listServices")
	}

	var dc = GetConfig().DatasourcesConfig
	for _, service := range services {
		typ := DatasourceTypeNone

		// by annotation
		for key, value := range service.Annotations {
			if key != dc.Discovery.AnnotationKey {
				continue
			}
			if value == string(DatasourceTypePrometheus) {
				typ = DatasourceTypePrometheus
				break
			}
			if value == string(DatasourceTypeLethe) {
				typ = DatasourceTypeLethe
				break
			}
		}
		// by name prometheus
		if typ == DatasourceTypeNone && dc.Discovery.ByNamePrometheus && service.Name == "prometheus" {
			typ = DatasourceTypePrometheus
		}
		// by name lethe
		if typ == DatasourceTypeNone && dc.Discovery.ByNameLethe && service.Name == "lethe" {
			typ = DatasourceTypeLethe
		}
		// not matched
		if typ == DatasourceTypeNone {
			continue
		}

		// isDefault
		isDefault := false
		if service.Namespace == dc.Discovery.DefaultNamespace {
			isDefault = true
		}

		// portNumber
		var portNumber int32 = 0
		for _, port := range service.Spec.Ports {
			if port.Name == "http" {
				portNumber = port.Port
				break
			}
		}
		if portNumber == 0 {
			portNumber = service.Spec.Ports[0].Port
		}

		datasources = append(datasources, Datasource{
			Name:         fmt.Sprintf("%s.%s", service.Name, service.Namespace),
			Type:         typ,
			URL:          fmt.Sprintf("http://%s.%s:%d", service.Name, service.Namespace, portNumber),
			IsDiscovered: true,
			IsDefault:    isDefault,
		})
	}
	return datasources, nil
}

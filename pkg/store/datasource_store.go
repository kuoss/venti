package store

import (
	"context"
	"fmt"
	"log"

	"github.com/kuoss/venti/pkg/configuration"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// DatasourceStore
type DatasourceStore struct {
	config      *configuration.DatasourcesConfig
	datasources []*configuration.Datasource
}

// NewDatasourceStore return *DatasourceStore after service discovery (with k8s service)
func NewDatasourceStore(cfg *configuration.DatasourcesConfig) (*DatasourceStore, error) {
	s := &DatasourceStore{cfg, nil}
	s.load()
	return s, nil
}

func (s *DatasourceStore) load() {
	// load from config
	s.datasources = s.config.Datasources

	// load from discovery
	if s.config.Discovery.Enabled {
		s.datasources = append(s.datasources, s.discoverDatasources()...)
	}
	// set main datasources
	s.setMainDatasources()
}

func (s *DatasourceStore) discoverDatasources() []*configuration.Datasource {
	services, err := listServices()
	if err != nil {
		log.Fatalf("error on listServices: %s", err)
		return []*configuration.Datasource{}
	}
	return s.getDatasourcesFromServices(services)
}

func listServices() ([]v1.Service, error) {
	config, err := rest.InClusterConfig()
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

func (s *DatasourceStore) getDatasourcesFromServices(services []v1.Service) []*configuration.Datasource {
	var datasources []*configuration.Datasource

	for _, service := range services {
		datasourceType := configuration.DatasourceTypeNone

		// by annotation
		for key, value := range service.Annotations {
			if key != s.config.Discovery.AnnotationKey {
				continue
			}
			if value == string(configuration.DatasourceTypePrometheus) {
				datasourceType = configuration.DatasourceTypePrometheus
				break
			}
			if value == string(configuration.DatasourceTypeLethe) {
				datasourceType = configuration.DatasourceTypeLethe
				break
			}
		}
		// by name prometheus
		if datasourceType == configuration.DatasourceTypeNone && s.config.Discovery.ByNamePrometheus && service.Name == "prometheus" {
			datasourceType = configuration.DatasourceTypePrometheus
		}
		// by name lethe
		if datasourceType == configuration.DatasourceTypeNone && s.config.Discovery.ByNameLethe && service.Name == "lethe" {
			datasourceType = configuration.DatasourceTypeLethe
		}
		// not matched
		if datasourceType == configuration.DatasourceTypeNone {
			continue
		}

		// isMain
		isMain := false
		if service.Namespace == s.config.Discovery.MainNamespace {
			isMain = true
		}

		portNumber := s.getPortNumberFromService(service)

		datasources = append(datasources, &configuration.Datasource{
			Name:         fmt.Sprintf("%s.%s", service.Name, service.Namespace),
			Type:         datasourceType,
			URL:          fmt.Sprintf("http://%s.%s:%d", service.Name, service.Namespace, portNumber),
			IsDiscovered: true,
			IsMain:       isMain,
		})
	}
	return datasources
}

// return port number within "http" named port. if not exist return service's first port number
func (s *DatasourceStore) getPortNumberFromService(service v1.Service) int32 {

	for _, port := range service.Spec.Ports {
		if port.Name == "http" {
			return port.Port
		}
	}
	return service.Spec.Ports[0].Port
}

// ensure that there is one main datasource for each type
func (s *DatasourceStore) setMainDatasources() {

	existsMainPrometheus := false
	existsMainLethe := false

	for _, ds := range s.datasources {
		if !ds.IsMain {
			continue
		}
		switch ds.Type {
		case configuration.DatasourceTypePrometheus:
			existsMainPrometheus = true
			continue
		case configuration.DatasourceTypeLethe:
			existsMainLethe = true
			continue
		}
	}
	if !existsMainPrometheus {
		for _, ds := range s.datasources {
			if ds.Type == configuration.DatasourceTypePrometheus {
				ds.IsMain = true
				break
			}
		}
	}
	if !existsMainLethe {
		for _, ds := range s.datasources {
			if ds.Type == configuration.DatasourceTypeLethe {
				ds.IsMain = true
				break
			}
		}
	}
}


func (s *DatasourceStore) GetDatasources() []*configuration.Datasource {
	return s.datasources
}

func (s *DatasourceStore) GetMainDatasourceWithType(typ configuration.DatasourceType) (*configuration.Datasource, error) {
	for _, ds := range s.datasources {
		if ds.Type == typ && ds.IsMain {
			return ds, nil
		}
	}
	return nil, fmt.Errorf("datasource of type %s not found", typ)
}

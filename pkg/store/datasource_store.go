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
	datasources []configuration.Datasource
}

// NewDatasourceStore return *DatasourceStore after service discovery (with k8s service)
func NewDatasourceStore(cfg *configuration.DatasourcesConfig) (*DatasourceStore, error) {
	store := &DatasourceStore{cfg, nil}
	err := store.load()
	if err != nil {
		return nil, fmt.Errorf("error on load: %w", err)
	}
	return store, nil
}

func (s *DatasourceStore) load() error {
	// load from config
	for _, datasource := range s.config.Datasources {
		s.datasources = append(s.datasources, *datasource)
	}

	// load from discovery
	if s.config.Discovery.Enabled {
		discoveredDatasources, err := s.discoverDatasources()
		if err != nil {
			log.Fatalf("error on discoverDatasources: %s", err)
		}
		s.datasources = append(s.datasources, discoveredDatasources...)
	}
	// set main datasources
	s.setMainDatasources()

	return nil
}

func (s *DatasourceStore) discoverDatasources() ([]configuration.Datasource, error) {
	services, err := listServices()
	if err != nil {
		return nil, fmt.Errorf("error on listServices: %s", err)
	}
	return s.getDatasourcesFromServices(services), nil
}

func listServices() ([]v1.Service, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot InClusterConfig: %w", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("cannot NewForConfig: %w", err)
	}

	services, err := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("cannot ListServices: %w", err)
	}
	return services.Items, nil
}

// get datasources from k8s services
// recognize as a datasource by annotation or name
func (s *DatasourceStore) getDatasourcesFromServices(services []v1.Service) []configuration.Datasource {
	var datasources []configuration.Datasource

	for _, service := range services {
		datasourceType := configuration.DatasourceTypeNone

		// recognize as a datasource by annotation of k8s service
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
		// recognize as a datasource by name 'prometheus'
		if datasourceType == configuration.DatasourceTypeNone && s.config.Discovery.ByNamePrometheus && service.Name == "prometheus" {
			datasourceType = configuration.DatasourceTypePrometheus
		}

		// recognize as a datasource by name 'lethe'
		if datasourceType == configuration.DatasourceTypeNone && s.config.Discovery.ByNameLethe && service.Name == "lethe" {
			datasourceType = configuration.DatasourceTypeLethe
		}
		// the service is not a datasource
		if datasourceType == configuration.DatasourceTypeNone {
			continue
		}

		// recognize as a main datasource by namespace
		isMain := false
		if service.Namespace == s.config.Discovery.MainNamespace {
			isMain = true
		}

		// get port number of datasource from k8s service
		portNumber := s.getPortNumberFromService(service)

		// append to datasources
		datasources = append(datasources, configuration.Datasource{
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

	// fallback for main prometheus datasource
	// If there is no main prometheus, the first prometheus will be a main prometheus.
	if !existsMainPrometheus {
		for _, ds := range s.datasources {
			if ds.Type == configuration.DatasourceTypePrometheus {
				ds.IsMain = true
				break
			}
		}
	}

	// fallback for main lethe datasource
	// If there is no main lethe, the first lethe will be a main lethe.
	if !existsMainLethe {
		for _, ds := range s.datasources {
			if ds.Type == configuration.DatasourceTypeLethe {
				ds.IsMain = true
				break
			}
		}
	}
}

func (s *DatasourceStore) GetDatasources() []configuration.Datasource {
	return s.datasources
}

func (s *DatasourceStore) GetMainDatasourceWithType(typ configuration.DatasourceType) (configuration.Datasource, error) {
	for _, ds := range s.datasources {
		if ds.Type == typ && ds.IsMain {
			return ds, nil
		}
	}
	return configuration.Datasource{}, fmt.Errorf("datasource of type %s not found", typ)
}

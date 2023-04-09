package store

import (
	"context"
	"fmt"
	"log"

	"github.com/kuoss/venti/pkg/model"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// DatasourceStore
type DatasourceStore struct {
	config      *model.DatasourcesConfig
	datasources []model.Datasource
}

// NewDatasourceStore return *DatasourceStore after service discovery (with k8s service)
func NewDatasourceStore(cfg *model.DatasourcesConfig) (*DatasourceStore, error) {
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

		clusterCfg, err := rest.InClusterConfig()
		if err != nil {
			return fmt.Errorf("cannot InClusterConfig: %w", err)
		}
		clientset, err := kubernetes.NewForConfig(clusterCfg)
		if err != nil {
			return fmt.Errorf("cannot NewForConfig: %w", err)
		}

		discoveredDatasources, err := s.discoverDatasources(clientset)
		if err != nil {
			log.Fatalf("error on discoverDatasources: %s", err)
		}
		s.datasources = append(s.datasources, discoveredDatasources...)
	}
	// set main datasources
	s.setMainDatasources()

	return nil
}

func (s *DatasourceStore) discoverDatasources(clientset kubernetes.Interface) ([]model.Datasource, error) {

	services, err := clientset.CoreV1().Services("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("cannot ListServices: %w", err)
	}
	return s.getDatasourcesFromServices(services.Items), nil
}

// get datasources from k8s servicesWithoutAnnotation
// recognize as a datasource by annotation or name
func (s *DatasourceStore) getDatasourcesFromServices(services []v1.Service) []model.Datasource {
	var datasources []model.Datasource

	for _, service := range services {
		datasourceType := getDatasourceTypeByConfig(service, s.config.Discovery)

		// the service is not a datasource
		if datasourceType == model.DatasourceTypeNone {
			continue
		}

		// recognize as a main datasource by namespace
		isMain := false
		if service.Namespace == s.config.Discovery.MainNamespace {
			isMain = true
		}

		// get port number of datasource from k8s service
		portNumber := getPortNumberFromService(service)

		// append to datasources
		datasources = append(datasources, model.Datasource{
			Name:         fmt.Sprintf("%s.%s", service.Name, service.Namespace),
			Type:         datasourceType,
			URL:          fmt.Sprintf("http://%s.%s:%d", service.Name, service.Namespace, portNumber),
			IsDiscovered: true,
			IsMain:       isMain,
		})
	}
	return datasources
}

// getDatasourceTypeByConfig return DatasourceType.
// 1. If configured within config.Discovery.ByNamePrometheus or config.Discovery.ByNameLethe return if service has matched name.
// 2. If configured within config.Discovery.AnnotationKey matched with service's annotation key and also value is
// one of promethe or lethe.
func getDatasourceTypeByConfig(service v1.Service, cfg model.Discovery) model.DatasourceType {

	// recognize as a datasource by name 'prometheus'
	if cfg.ByNamePrometheus && service.Name == "prometheus" {
		return model.DatasourceTypePrometheus
	}
	// recognize as a datasource by name 'lethe'
	if cfg.ByNameLethe && service.Name == "lethe" {
		return model.DatasourceTypeLethe
	}

	// recognize as a datasource by annotation of k8s service
	for key, value := range service.Annotations {
		if key != cfg.AnnotationKey {
			continue
		}
		if value == string(model.DatasourceTypePrometheus) {
			return model.DatasourceTypePrometheus
		}
		if value == string(model.DatasourceTypeLethe) {
			return model.DatasourceTypeLethe
		}
	}

	return model.DatasourceTypeNone
}

// return port number within "http" named port. if not exist return service's first port number
func getPortNumberFromService(service v1.Service) int32 {

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
		case model.DatasourceTypePrometheus:
			existsMainPrometheus = true
			continue
		case model.DatasourceTypeLethe:
			existsMainLethe = true
			continue
		}
	}

	// fallback for main prometheus datasource
	// If there is no main prometheus, the first prometheus will be a main prometheus.
	if !existsMainPrometheus {
		for _, ds := range s.datasources {
			if ds.Type == model.DatasourceTypePrometheus {
				ds.IsMain = true
				break
			}
		}
	}

	// fallback for main lethe datasource
	// If there is no main lethe, the first lethe will be a main lethe.
	if !existsMainLethe {
		for _, ds := range s.datasources {
			if ds.Type == model.DatasourceTypeLethe {
				ds.IsMain = true
				break
			}
		}
	}
}

func (s *DatasourceStore) GetDatasources() []model.Datasource {
	return s.datasources
}

func (s *DatasourceStore) GetMainDatasourceWithType(typ model.DatasourceType) (model.Datasource, error) {
	for _, ds := range s.datasources {
		if ds.Type == typ && ds.IsMain {
			return ds, nil
		}
	}
	return model.Datasource{}, fmt.Errorf("datasource of type %s not found", typ)
}

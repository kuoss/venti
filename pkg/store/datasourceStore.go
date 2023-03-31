package store

import (
	"context"
	"fmt"
	"github.com/kuoss/venti/pkg/configuration"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// TODO can't fully understand initializing Datasource

// DatasourceStore
type DatasourceStore struct {
	*configuration.Config
}

func NewDatasourceStore(config *configuration.Config) *DatasourceStore {
	dss := &DatasourceStore{config}
	dss.setDefaultDatasources()
	return dss
}

// TODO do we need this?
func (d *DatasourceStore) GetDatasources() []*configuration.Datasource {
	return d.DatasourcesConfig.Datasources
}

func (d *DatasourceStore) GetDefaultDatasource(typ configuration.DatasourceType) (*configuration.Datasource, error) {
	for _, ds := range d.DatasourcesConfig.Datasources {
		if ds.Type == typ && ds.IsDefault {
			return ds, nil
		}
	}
	return nil, fmt.Errorf("datasource of type %s not found", typ)
}

func (d *DatasourceStore) loadDatasources() error {
	if d.DatasourcesConfig.Discovery.Enabled {
		services, err := d.listServices()
		if err != nil {
			return fmt.Errorf("cannot listServices: %w", err)
		}

		discoveredDatasources := d.getDatasourcesFromServices(services)
		d.DatasourcesConfig.Datasources = append(d.DatasourcesConfig.Datasources, discoveredDatasources...)
	}
	return nil
}

func (d *DatasourceStore) listServices() ([]v1.Service, error) {
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

func (d *DatasourceStore) getDatasourcesFromServices(services []v1.Service) []*configuration.Datasource {
	var datasources []*configuration.Datasource

	for _, service := range services {
		datasourceType := configuration.DatasourceTypeNone

		// by annotation
		for key, value := range service.Annotations {
			if key != d.DatasourcesConfig.Discovery.AnnotationKey {
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
		if datasourceType == configuration.DatasourceTypeNone && d.DatasourcesConfig.Discovery.ByNamePrometheus && service.Name == "prometheus" {
			datasourceType = configuration.DatasourceTypePrometheus
		}
		// by name lethe
		if datasourceType == configuration.DatasourceTypeNone && d.DatasourcesConfig.Discovery.ByNameLethe && service.Name == "lethe" {
			datasourceType = configuration.DatasourceTypeLethe
		}
		// not matched
		if datasourceType == configuration.DatasourceTypeNone {
			continue
		}

		// isDefault
		isDefault := false
		if service.Namespace == d.DatasourcesConfig.Discovery.DefaultNamespace {
			isDefault = true
		}

		// portNumber
		//	var portNumber int32 = 0

		portNumber := getPortNumberFromService(service)

		datasources = append(datasources, &configuration.Datasource{
			Name:         fmt.Sprintf("%s.%s", service.Name, service.Namespace),
			Type:         datasourceType,
			URL:          fmt.Sprintf("http://%s.%s:%d", service.Name, service.Namespace, portNumber),
			IsDiscovered: true,
			IsDefault:    isDefault,
		})
	}
	return datasources
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

func (d *DatasourceStore) setDefaultDatasources() {
	// initialize with false
	var (
		existsDefaultPrometheus, existsDefaultLethe bool
	)

	for _, ds := range d.DatasourcesConfig.Datasources {
		if !ds.IsDefault {
			continue
		}
		switch ds.Type {
		case configuration.DatasourceTypePrometheus:
			existsDefaultPrometheus = true
			continue
		case configuration.DatasourceTypeLethe:
			existsDefaultLethe = true
			continue
		}
	}
	if !existsDefaultPrometheus {
		for _, ds := range d.DatasourcesConfig.Datasources {
			if ds.Type == configuration.DatasourceTypePrometheus {
				ds.IsDefault = true
				break
			}
		}
	}
	if !existsDefaultLethe {
		for _, ds := range d.DatasourcesConfig.Datasources {
			if ds.Type == configuration.DatasourceTypeLethe {
				ds.IsDefault = true
				break
			}
		}
	}
}

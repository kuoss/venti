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
	config     *configuration.DatasourcesConfig
	discovered []*configuration.Datasource // discovered
}

// NewDatasourceStore return *DatasourceStore after service discovery (with k8s service)
func NewDatasourceStore(config *configuration.DatasourcesConfig) (*DatasourceStore, error) {
	dss := &DatasourceStore{config, nil}
	err := dss.discovery()
	if err != nil {
		return nil, fmt.Errorf("failed to service discovery: %w", err)
	}
	dss.setDefaultDatasources()
	return dss, nil
}

// TODO do we need this?

func (d *DatasourceStore) GetDatasources() []*configuration.Datasource {
	return d.discovered
}

func (d *DatasourceStore) GetDatasourceWithType(t configuration.DatasourceType) *configuration.Datasource {
	for _, ds := range d.discovered {
		if ds.Type == t {
			return ds
		}
	}
	return nil
}

func (d *DatasourceStore) GetDefaultDatasource(typ configuration.DatasourceType) (*configuration.Datasource, error) {
	for _, ds := range d.discovered {
		if ds.Type == typ && ds.IsDefault {
			return ds, nil
		}
	}
	return nil, fmt.Errorf("datasource of type %s not found", typ)
}

func (d *DatasourceStore) discovery() error {
	if d.config.Discovery.Enabled {
		services, err := listServices()
		if err != nil {
			return fmt.Errorf("cannot listServices: %w", err)
		}

		discovered := d.getDatasourcesFromServices(services)
		d.discovered = append(d.discovered, discovered...)
	}
	return nil
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

func (d *DatasourceStore) getDatasourcesFromServices(services []v1.Service) []*configuration.Datasource {
	var discovered []*configuration.Datasource

	for _, service := range services {
		datasourceType := configuration.DatasourceTypeNone

		// by annotation
		for key, value := range service.Annotations {
			if key != d.config.Discovery.AnnotationKey {
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
		if datasourceType == configuration.DatasourceTypeNone && d.config.Discovery.ByNamePrometheus && service.Name == "prometheus" {
			datasourceType = configuration.DatasourceTypePrometheus
		}
		// by name lethe
		if datasourceType == configuration.DatasourceTypeNone && d.config.Discovery.ByNameLethe && service.Name == "lethe" {
			datasourceType = configuration.DatasourceTypeLethe
		}
		// not matched
		if datasourceType == configuration.DatasourceTypeNone {
			continue
		}

		// isDefault
		isDefault := false
		if service.Namespace == d.config.Discovery.DefaultNamespace {
			isDefault = true
		}

		// portNumber
		//	var portNumber int32 = 0

		portNumber := getPortNumberFromService(service)

		discovered = append(discovered, &configuration.Datasource{
			Name:         fmt.Sprintf("%s.%s", service.Name, service.Namespace),
			Type:         datasourceType,
			URL:          fmt.Sprintf("http://%s.%s:%d", service.Name, service.Namespace, portNumber),
			IsDiscovered: true,
			IsDefault:    isDefault,
		})
	}
	return discovered
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

// todo : what is this for?
func (d *DatasourceStore) setDefaultDatasources() {

	// initialize with false
	var (
		existsDefaultPrometheus, existsDefaultLethe bool
	)

	for _, ds := range d.discovered {
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
		for _, ds := range d.discovered {
			if ds.Type == configuration.DatasourceTypePrometheus {
				ds.IsDefault = true
				break
			}
		}
	}
	if !existsDefaultLethe {
		for _, ds := range d.discovered {
			if ds.Type == configuration.DatasourceTypeLethe {
				ds.IsDefault = true
				break
			}
		}
	}
}

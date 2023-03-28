package server

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
)

type DatasourceStore struct {
	datasourcesConfig DatasourcesConfig
	datasources       []Datasource
}

func NewDatasourceStore(datasourcesConfig DatasourcesConfig) *DatasourceStore {
	dss := &DatasourceStore{datasourcesConfig: datasourcesConfig}
	dss.loadDatasources()
	dss.setDefaultDatasources()
	return dss
}

func (d *DatasourceStore) GetDatasources() []Datasource {
	return d.datasources
}

func (d *DatasourceStore) GetDefaultDatasource(typ DatasourceType) (Datasource, error) {
	for _, ds := range d.datasources {
		if ds.Type == typ && ds.IsDefault {
			return ds, nil
		}
	}
	return Datasource{}, fmt.Errorf("datasource of type %s not found", typ)
}

func (d *DatasourceStore) loadDatasources() {
	d.datasources = d.datasourcesConfig.Datasources
	if d.datasourcesConfig.Discovery.Enabled {
		services, err := d.discoverServices()
		if err != nil {
			log.Printf("cannot discoverServices: %s", err)
		} else {
			discoveredDatasources := d.getDatasourcesFromServices(services)
			d.datasources = append(d.datasources, discoveredDatasources...)
		}
	}
}

func (d *DatasourceStore) discoverServices() ([]v1.Service, error) {
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

func (d *DatasourceStore) getDatasourcesFromServices(services []v1.Service) []Datasource {
	var datasources = []Datasource{}

	for _, service := range services {
		typ := DatasourceTypeNone

		// by annotation
		for key, value := range service.Annotations {
			if key != d.datasourcesConfig.Discovery.AnnotationKey {
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
		if typ == DatasourceTypeNone && d.datasourcesConfig.Discovery.ByNamePrometheus && service.Name == "prometheus" {
			typ = DatasourceTypePrometheus
		}
		// by name lethe
		if typ == DatasourceTypeNone && d.datasourcesConfig.Discovery.ByNameLethe && service.Name == "lethe" {
			typ = DatasourceTypeLethe
		}
		// not matched
		if typ == DatasourceTypeNone {
			continue
		}

		// isDefault
		isDefault := false
		if service.Namespace == d.datasourcesConfig.Discovery.DefaultNamespace {
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
	return datasources
}

func (d *DatasourceStore) setDefaultDatasources() {
	existsDefaultPrometheus := false
	existsDefaultLethe := false
	for _, ds := range d.datasources {
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
		for i, ds := range d.datasources {
			if ds.Type == DatasourceTypePrometheus {
				d.datasources[i].IsDefault = true
				break
			}
		}
	}
	if !existsDefaultLethe {
		for i, ds := range d.datasources {
			if ds.Type == DatasourceTypeLethe {
				d.datasources[i].IsDefault = true
				break
			}
		}
	}
}

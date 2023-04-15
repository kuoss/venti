package kubernetes

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

type k8sStore struct {
	client kubernetes.Interface
}

func NewK8sStore() (*k8sStore, error) {
	clusterCfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot InClusterConfig: %w", err)
	}
	clientset, err := kubernetes.NewForConfig(clusterCfg)
	if err != nil {
		return nil, fmt.Errorf("cannot NewForConfig: %w", err)
	}
	return &k8sStore{client: clientset}, nil
}

func (s *k8sStore) Do(discovery model.Discovery) ([]model.Datasource, error) {

	services, err := s.client.CoreV1().Services("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("cannot ListServices: %w", err)
	}
	return s.getDatasourcesFromServices(services.Items, discovery), nil
}

func (s *k8sStore) getDatasourcesFromServices(services []v1.Service, discovery model.Discovery) []model.Datasource {
	var datasources []model.Datasource

	for _, service := range services {
		datasourceType := getDatasourceTypeByConfig(service, discovery)

		// the service is not a datasource
		if datasourceType == model.DatasourceTypeNone {
			continue
		}

		// recognize as a main datasource by namespace
		isMain := false
		if service.Namespace == discovery.MainNamespace {
			isMain = true
		}

		// get port number of datasource from k8s service
		portNumber, err := getPortNumberFromService(service)
		if err != nil {
			log.Printf("extract port number from service failed. %s", err)
			continue
		}
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
func getPortNumberFromService(service v1.Service) (int32, error) {
	if len(service.Spec.Ports) < 1 {
		return 0, fmt.Errorf("service %s/%s have any port", service.Namespace, service.Name)
	}

	for _, port := range service.Spec.Ports {
		if port.Name == "http" {
			return port.Port, nil
		}
	}
	return service.Spec.Ports[0].Port, nil
}

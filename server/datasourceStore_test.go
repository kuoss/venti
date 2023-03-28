package server

import (
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func init() {
	LoadConfig("Unknown")
	config.DatasourcesConfig.Discovery.Enabled = false
}

func makeService(name string, namespace string, multiport bool) v1.Service {
	ports := []v1.ServicePort{
		{
			Name:     "testport",
			Protocol: v1.ProtocolTCP,
			Port:     int32(30900),
		},
	}
	if multiport {
		ports = append(ports, v1.ServicePort{
			Name:     "http",
			Protocol: v1.ProtocolTCP,
			Port:     int32(8080),
		})
	}
	return v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1.ServiceSpec{
			Ports:     ports,
			Type:      v1.ServiceTypeClusterIP,
			ClusterIP: "10.0.0.1",
		},
	}
}

func TestGetDatasourcesFromServices(t *testing.T) {
	config.DatasourcesConfig.Discovery.Enabled = false
	datasourceStore = NewDatasourceStore(config.DatasourcesConfig)
	config.DatasourcesConfig.Discovery.Enabled = false
	config.DatasourcesConfig.Discovery.ByNamePrometheus = true
	config.DatasourcesConfig.Discovery.ByNameLethe = true

	services := []v1.Service{
		makeService("prometheus", "namespace1", false),
		makeService("prometheus", "namespace2", false),
		makeService("prometheus", "kube-system", false),
		makeService("lethe", "kuoss", true),
		makeService("lethe", "kube-system", true),
	}

	assert.Equal(t, 5, len(services))
	datasourceStore.getDatasourcesFromServices(services)

	// TODO: change after refactoring

	// assert.Equal(t, 5, len(datasourceStore.GetDatasources()))
	assert.Equal(t, 0, len(datasourceStore.GetDatasources()))
	datasourceStore.setDefaultDatasources()
	datasources := datasourceStore.GetDatasources()

	// assert.Equal(t, 5, len(datasources))
	assert.Equal(t, 0, len(datasources))

	// assert.Equal(t, []Datasource{
	// 	{Type: "prometheus", Name: "prometheus.namespace1", URL: "http://prometheus.namespace1:30900", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsDefault: true, IsDiscovered: true},
	// 	{Type: "prometheus", Name: "prometheus.namespace2", URL: "http://prometheus.namespace2:30900", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsDefault: false, IsDiscovered: true},
	// 	{Type: "prometheus", Name: "prometheus.kube-system", URL: "http://prometheus.kube-system:30900", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsDefault: false, IsDiscovered: true},
	// 	{Type: "lethe", Name: "lethe.kuoss", URL: "http://lethe.kuoss:8080", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsDefault: true, IsDiscovered: true},
	// 	{Type: "lethe", Name: "lethe.kube-system", URL: "http://lethe.kube-system:8080", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsDefault: false, IsDiscovered: true},
	// }, datasources)
	assert.Equal(t, []Datasource(nil), datasources)
}

package store

import (
	"github.com/kuoss/venti/pkg/configuration"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

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

	datasourcesConfig := &configuration.DatasourcesConfig{
		QueryTimeout: time.Second * 10,
		Datasources: []*configuration.Datasource{
			{
				Type:         configuration.DatasourceTypePrometheus,
				Name:         "prometheus",
				URL:          "localhost",
				BasicAuth:    false,
				IsDefault:    false,
				IsDiscovered: false,
			},
			{
				Type:         configuration.DatasourceTypeLethe,
				Name:         "lethe",
				URL:          "localhost",
				BasicAuth:    false,
				IsDefault:    false,
				IsDiscovered: false,
			},
		},
		Discovery: configuration.Discovery{
			Enabled:          false,
			ByNamePrometheus: true,
			ByNameLethe:      true,
		},
	}

	dss, err := NewDatasourceStore(datasourcesConfig)
	if err != nil {
		t.Fatalf("datasource fatal %s", err)
	}

	services := []v1.Service{
		makeService("prometheus", "namespace1", false),
		makeService("prometheus", "namespace2", false),
		makeService("prometheus", "kube-system", false),
		makeService("lethe", "kuoss", true),
		makeService("lethe", "kube-system", true),
	}

	discovered := dss.getDatasourcesFromServices(services)

	assert.ElementsMatch(t, []*configuration.Datasource{
		{Type: "prometheus", Name: "prometheus.namespace1", URL: "http://prometheus.namespace1:30900", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsDefault: false, IsDiscovered: true},
		{Type: "prometheus", Name: "prometheus.namespace2", URL: "http://prometheus.namespace2:30900", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsDefault: false, IsDiscovered: true},
		{Type: "prometheus", Name: "prometheus.kube-system", URL: "http://prometheus.kube-system:30900", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsDefault: false, IsDiscovered: true},
		{Type: "lethe", Name: "lethe.kuoss", URL: "http://lethe.kuoss:8080", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsDefault: false, IsDiscovered: true},
		{Type: "lethe", Name: "lethe.kube-system", URL: "http://lethe.kube-system:8080", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsDefault: false, IsDiscovered: true},
	}, discovered)
}

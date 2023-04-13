package kubernetes

import (
	"github.com/kuoss/venti/pkg/model"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func makeService(name string, namespace string, multiport bool, annotation map[string]string) runtime.Object {
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

	return runtime.Object(&v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: annotation,
		},
		Spec: v1.ServiceSpec{
			Ports:     ports,
			Type:      v1.ServiceTypeClusterIP,
			ClusterIP: "10.0.0.1",
		},
	})
}

var servicesWithoutAnnotation = []runtime.Object{
	makeService("prometheus", "namespace1", false, nil),
	makeService("prometheus", "namespace2", false, nil),
	makeService("prometheus", "kube-system", false, nil),
	makeService("lethe", "kuoss", true, nil),
	makeService("lethe", "kube-system", true, nil),
}

var servicesWithAnnotation = []runtime.Object{
	makeService("prometheus", "namespace1", false, map[string]string{
		"kuoss.org/datasource-type": "prometheus",
	}),
	makeService("prometheus", "namespace2", false, map[string]string{
		"kuoss.org/datasource-type": "prometheus",
	}),
	makeService("prometheus", "kube-system", false, map[string]string{
		"kuoss.org/datasource-type": "prometheus",
	}),
	makeService("lethe", "kuoss", true, map[string]string{
		"kuoss.org/datasource-type": "lethe",
	}),
	makeService("lethe", "kube-system", true, map[string]string{
		"kuoss.org/datasource-type": "lethe",
	}),
}

func TestDoDiscoveryWithoutAnnotationKey(t *testing.T) {
	want := []model.Datasource{
		{
			Type:         "lethe",
			Name:         "lethe.kube-system",
			URL:          "http://lethe.kube-system:8080",
			IsMain:       false,
			IsDiscovered: true,
		},
		{
			Type:         "prometheus",
			Name:         "prometheus.kube-system",
			URL:          "http://prometheus.kube-system:30900",
			IsMain:       false,
			IsDiscovered: true,
		},
		{
			Type:         "lethe",
			Name:         "lethe.kuoss",
			URL:          "http://lethe.kuoss:8080",
			IsMain:       false,
			IsDiscovered: true,
		},
		{
			Type:         "prometheus",
			Name:         "prometheus.namespace1",
			URL:          "http://prometheus.namespace1:30900",
			IsMain:       false,
			IsDiscovered: true,
		},
		{
			Type:         "prometheus",
			Name:         "prometheus.namespace2",
			URL:          "http://prometheus.namespace2:30900",
			IsMain:       false,
			IsDiscovered: true,
		}}

	k8sStore := &k8sStore{fake.NewSimpleClientset(servicesWithoutAnnotation...)}
	discovered, err := k8sStore.Do(model.Discovery{
		Enabled:          true,
		ByNamePrometheus: true,
		ByNameLethe:      true,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.ElementsMatch(t, want, discovered)
}

func TestDoDiscoveryAnnotationKey(t *testing.T) {
	want := []model.Datasource{
		{
			Type:         "lethe",
			Name:         "lethe.kube-system",
			URL:          "http://lethe.kube-system:8080",
			IsMain:       false,
			IsDiscovered: true,
		},
		{
			Type:         "prometheus",
			Name:         "prometheus.kube-system",
			URL:          "http://prometheus.kube-system:30900",
			IsMain:       false,
			IsDiscovered: true,
		},
		{
			Type:         "lethe",
			Name:         "lethe.kuoss",
			URL:          "http://lethe.kuoss:8080",
			IsMain:       false,
			IsDiscovered: true,
		},
		{
			Type:         "prometheus",
			Name:         "prometheus.namespace1",
			URL:          "http://prometheus.namespace1:30900",
			IsMain:       false,
			IsDiscovered: true,
		},
		{
			Type:         "prometheus",
			Name:         "prometheus.namespace2",
			URL:          "http://prometheus.namespace2:30900",
			IsMain:       false,
			IsDiscovered: true,
		}}

	k8sStore := &k8sStore{fake.NewSimpleClientset(servicesWithAnnotation...)}
	discovered, err := k8sStore.Do(model.Discovery{
		Enabled:          true,
		AnnotationKey:    "kuoss.org/datasource-type",
		ByNamePrometheus: false,
		ByNameLethe:      false,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.ElementsMatch(t, want, discovered)
}

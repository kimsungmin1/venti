package store

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
	"time"

	"github.com/kuoss/venti/pkg/configuration"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func TestNewDatasourceStore(t *testing.T) {
	datasources := []configuration.Datasource{
		{Type: configuration.DatasourceTypePrometheus, Name: "Prometheus", URL: "http://prometheus:9090", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsMain: false, IsDiscovered: false},
		{Type: configuration.DatasourceTypeLethe, Name: "Lethe", URL: "http://lethe:3100", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsMain: false, IsDiscovered: false},
	}
	datasourcesPointer := []*configuration.Datasource{
		{Type: configuration.DatasourceTypePrometheus, Name: "Prometheus", URL: "http://prometheus:9090", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsMain: false, IsDiscovered: false},
		{Type: configuration.DatasourceTypeLethe, Name: "Lethe", URL: "http://lethe:3100", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsMain: false, IsDiscovered: false},
	}
	datasourcesConfig := &configuration.DatasourcesConfig{
		QueryTimeout: time.Second * 10,
		Datasources:  datasourcesPointer,
		Discovery: configuration.Discovery{
			Enabled:          false,
			ByNamePrometheus: true,
			ByNameLethe:      true,
		},
	}
	store, err := NewDatasourceStore(datasourcesConfig)
	assert.Nil(t, err)
	assert.Equal(t, store.config, datasourcesConfig)
	assert.ElementsMatch(t, store.datasources, datasources)
}

var servicesWithoutAnnotation = []runtime.Object{
	makeService("prometheus", "namespace1", false, nil),
	makeService("prometheus", "namespace2", false, nil),
	makeService("prometheus", "kube-system", false, nil),
	makeService("lethe", "kuoss", true, nil),
	makeService("lethe", "kube-system", true, nil),
}

func TestDiscoverDatasourcesByName(t *testing.T) {
	datasourcesConfig := &configuration.DatasourcesConfig{
		QueryTimeout: time.Second * 10,
		Datasources:  []*configuration.Datasource{},
		Discovery: configuration.Discovery{
			Enabled:          false,
			ByNamePrometheus: true,
			ByNameLethe:      true,
		},
	}
	want := []configuration.Datasource{
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
	store, _ := NewDatasourceStore(datasourcesConfig)
	got, err := store.discoverDatasources(fake.NewSimpleClientset(servicesWithoutAnnotation...))
	if err != nil {
		return
	}
	assert.ElementsMatch(t, want, got)
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

func TestDiscoverDatasourcesByAnnotationKey(t *testing.T) {
	datasourcesConfig := &configuration.DatasourcesConfig{
		QueryTimeout: time.Second * 10,
		Datasources:  []*configuration.Datasource{},
		Discovery: configuration.Discovery{
			Enabled:          false, // cheat
			AnnotationKey:    "kuoss.org/datasource-type",
			ByNamePrometheus: false,
			ByNameLethe:      false,
		},
	}
	want := []configuration.Datasource{
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
	store, _ := NewDatasourceStore(datasourcesConfig)
	got, err := store.discoverDatasources(fake.NewSimpleClientset(servicesWithAnnotation...))
	if err != nil {
		return
	}
	assert.ElementsMatch(t, got, want)
}

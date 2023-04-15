package store

import (
	"testing"
	"time"

	"github.com/kuoss/venti/pkg/store/discovery"

	"github.com/kuoss/venti/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestNewDatasourceStore(t *testing.T) {
	datasources := []model.Datasource{
		{Type: model.DatasourceTypePrometheus, Name: "Prometheus", URL: "http://prometheus:9090", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsMain: false, IsDiscovered: false},
		{Type: model.DatasourceTypeLethe, Name: "Lethe", URL: "http://lethe:3100", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsMain: false, IsDiscovered: false},
	}
	datasourcesPointer := []*model.Datasource{
		{Type: model.DatasourceTypePrometheus, Name: "Prometheus", URL: "http://prometheus:9090", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsMain: false, IsDiscovered: false},
		{Type: model.DatasourceTypeLethe, Name: "Lethe", URL: "http://lethe:3100", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsMain: false, IsDiscovered: false},
	}
	datasourcesConfig := &model.DatasourcesConfig{
		QueryTimeout: time.Second * 10,
		Datasources:  datasourcesPointer,
		Discovery: model.Discovery{
			Enabled:          false,
			ByNamePrometheus: true,
			ByNameLethe:      true,
		},
	}
	var defaultDiscoverer discovery.Discoverer
	store, err := NewDatasourceStore(datasourcesConfig, defaultDiscoverer)
	assert.Nil(t, err)
	assert.Equal(t, store.config, datasourcesConfig)
	assert.ElementsMatch(t, store.datasources, datasources)
}

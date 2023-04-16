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
	datasourcePointers := []*model.Datasource{
		{Type: model.DatasourceTypePrometheus, Name: "Prometheus", URL: "http://prometheus:9090", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsMain: false, IsDiscovered: false},
		{Type: model.DatasourceTypeLethe, Name: "Lethe", URL: "http://lethe:3100", BasicAuth: false, BasicAuthUser: "", BasicAuthPassword: "", IsMain: false, IsDiscovered: false},
	}
	datasourceConfig := &model.DatasourceConfig{
		QueryTimeout: time.Second * 10,
		Datasources:  datasourcePointers,
		Discovery: model.Discovery{
			Enabled:          false,
			ByNamePrometheus: true,
			ByNameLethe:      true,
		},
	}
	var defaultDiscoverer discovery.Discoverer
	store, err := NewDatasourceStore(datasourceConfig, defaultDiscoverer)
	assert.Nil(t, err)
	assert.Equal(t, store.config, datasourceConfig)
	assert.ElementsMatch(t, store.datasources, datasources)
}

package store

import (
	"os"
	"testing"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/venti/pkg/model"
	"github.com/stretchr/testify/assert"
)

func init() {
	err := os.Chdir("../..")
	if err != nil {
		panic(err)
	}
	logger.SetCallerSkip(9)
	logger.SetLevel(logger.DebugLevel)
	logger.Infof("init")
}

func TestNewStores(t *testing.T) {
	datasourceConfig := model.DatasourceConfig{
		Datasources: []model.Datasource{
			{Name: "mainPrometheus", Type: model.DatasourceTypePrometheus, URL: "http://prometheus:9090", IsMain: true},
			{Name: "subPrometheus1", Type: model.DatasourceTypePrometheus, URL: "http://prometheus1:9090", IsMain: false},
			{Name: "subPrometheus2", Type: model.DatasourceTypePrometheus, URL: "http://prometheus2:9090", IsMain: false},
			{Name: "mainLethe", Type: model.DatasourceTypeLethe, URL: "http://lethe:3100", IsMain: true},
			{Name: "subLethe1", Type: model.DatasourceTypeLethe, URL: "http://lethe1:3100", IsMain: false},
			{Name: "subLethe2", Type: model.DatasourceTypeLethe, URL: "http://lethe2:3100", IsMain: false},
		},
		Discovery: model.Discovery{
			Enabled:          false,
			ByNamePrometheus: true,
			ByNameLethe:      true,
		},
	}
	got, err := NewStores(&model.Config{DatasourceConfig: datasourceConfig})
	assert.NoError(t, err)
	assert.NotEmpty(t, got)
	assert.NotEmpty(t, got.AlertRuleStore)
	assert.NotEmpty(t, got.AlertingStore)
	assert.NotEmpty(t, got.DashboardStore)
	assert.NotEmpty(t, got.DatasourceStore)
	assert.NotEmpty(t, got.RemoteStore)
	assert.NotEmpty(t, got.StatusStore)
	assert.NotEmpty(t, got.UserStore)
}

func TestNewStoresError(t *testing.T) {
	got, err := NewStores(&model.Config{})
	assert.NoError(t, err)
	assert.NotEmpty(t, got)
	assert.NotEmpty(t, got.AlertRuleStore)
	assert.NotEmpty(t, got.AlertingStore)
	assert.NotEmpty(t, got.DashboardStore)
	assert.Empty(t, got.DatasourceStore)
	assert.NotEmpty(t, got.RemoteStore)
	assert.NotEmpty(t, got.StatusStore)
	assert.NotEmpty(t, got.UserStore)
}

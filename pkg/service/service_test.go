package service

import (
	"os"
	"testing"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/venti/pkg/model"
	"github.com/stretchr/testify/require"
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

func TestNewServices(t *testing.T) {
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
	got, err := NewServices(&model.Config{DatasourceConfig: datasourceConfig})
	require.NoError(t, err)
	require.NotEmpty(t, got)
	require.NotEmpty(t, got.AlertRuleService)
	require.NotEmpty(t, got.AlertingService)
	require.NotEmpty(t, got.DashboardService)
	require.NotEmpty(t, got.DatasourceService)
	require.NotEmpty(t, got.RemoteService)
	require.NotEmpty(t, got.StatusService)
	require.NotEmpty(t, got.UserService)
}

func TestNewServicesError(t *testing.T) {
	got, err := NewServices(&model.Config{})
	require.NoError(t, err)
	require.NotEmpty(t, got)
	require.NotEmpty(t, got.AlertRuleService)
	require.NotEmpty(t, got.AlertingService)
	require.NotEmpty(t, got.DashboardService)
	require.NotEmpty(t, got.DatasourceService)
	require.NotEmpty(t, got.RemoteService)
	require.NotEmpty(t, got.StatusService)
	require.NotEmpty(t, got.UserService)
}

package handler

import (
	"testing"

	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/service"
	"github.com/stretchr/testify/assert"
)

func TestLoadHandlers(t *testing.T) {

	cfg := &model.Config{
		AppInfo:          model.AppInfo{Version: "Unknown"},
		GlobalConfig:     model.GlobalConfig{},
		DatasourceConfig: model.DatasourceConfig{Datasources: []model.Datasource{{Type: model.DatasourceTypePrometheus, Name: "prometheus", IsMain: true}}},
		UserConfig:       model.UserConfig{},
		AlertingConfig:   model.AlertingConfig{},
	}
	services, err := service.NewServices(cfg)
	assert.NoError(t, err)
	assert.NotEmpty(t, services)

	handlers := loadHandlers(services)
	assert.NotEmpty(t, handlers)
	assert.NotEmpty(t, handlers.alertHandler)
	assert.NotEmpty(t, handlers.authHandler)
	assert.NotEmpty(t, handlers.dashboardHandler)
	assert.NotEmpty(t, handlers.datasourceHandler)
	assert.NotEmpty(t, handlers.remoteHandler)
	assert.NotEmpty(t, handlers.statusHandler)
}

package handler

import (
	"testing"

	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/store"
	"github.com/stretchr/testify/assert"
)

func TestLoadHandlers(t *testing.T) {

	cfg := &model.Config{
		Version:    "Unknown",
		UserConfig: model.UserConfig{},
		DatasourceConfig: model.DatasourceConfig{
			Datasources: []model.Datasource{
				{Type: model.DatasourceTypePrometheus, Name: "prometheus", IsMain: true},
			},
		},
	}
	stores, err := store.NewStores(cfg)
	assert.NoError(t, err)
	assert.NotEmpty(t, stores)

	handlers := loadHandlers(cfg, stores)
	assert.NotEmpty(t, handlers)
	assert.NotEmpty(t, handlers.alertHandler)
	assert.NotEmpty(t, handlers.authHandler)
	assert.NotEmpty(t, handlers.configHandler)
	assert.NotEmpty(t, handlers.dashboardHandler)
	assert.NotEmpty(t, handlers.datasourceHandler)
	assert.NotEmpty(t, handlers.remoteHandler)
	assert.NotEmpty(t, handlers.statusHandler)
}

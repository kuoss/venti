package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/mock"
	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/store"
	"github.com/kuoss/venti/pkg/store/discovery"
)

var (
	handlers          *Handlers
	datasource        *model.Datasource
	datasourcesConfig *model.DatasourcesConfig
	router            *gin.Engine
	stores            *store.Stores
)

func init() {
	_ = os.Chdir("../..")
	ts := mock.PrometheusServer()
	datasource = &model.Datasource{
		Type:   model.DatasourceTypePrometheus,
		Name:   "prometheus",
		URL:    ts.URL,
		IsMain: true,
	}
	datasourcesConfig = &model.DatasourcesConfig{
		Datasources:  []*model.Datasource{datasource},
		QueryTimeout: 30 * time.Second,
	}
	cfg := &model.Config{
		Version:           "Unknown",
		UserConfig:        model.UsersConfig{},
		DatasourcesConfig: datasourcesConfig,
	}
	var discoverer discovery.Discoverer
	datatsourceStore, _ := store.NewDatasourceStore(datasourcesConfig, discoverer)
	remoteStore := store.NewRemoteStore(&http.Client{}, datasourcesConfig.QueryTimeout)
	stores = &store.Stores{
		AlertRuleStore:  &store.AlertRuleStore{},
		DashboardStore:  &store.DashboardStore{},
		DatasourceStore: datatsourceStore,
		UserStore:       &store.UserStore{},
		RemoteStore:     remoteStore,
	}
	handlers = loadHandlers(cfg, stores)
	router = SetupRouter(cfg, stores)
}

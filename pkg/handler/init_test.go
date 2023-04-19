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
	handlers         *Handlers
	datasource       *model.Datasource
	datasourceConfig *model.DatasourceConfig
	router           *gin.Engine
	stores           *store.Stores
)

func init() {
	_ = os.Chdir("../..")
	ts := mock.Prometheus()
	datasource = &model.Datasource{
		Type:   model.DatasourceTypePrometheus,
		Name:   "prometheus",
		URL:    ts.URL,
		IsMain: true,
	}
	datasourceConfig = &model.DatasourceConfig{
		Datasources:  []*model.Datasource{datasource},
		QueryTimeout: 30 * time.Second,
	}
	cfg := &model.Config{
		Version:          "Unknown",
		UserConfig:       &model.UserConfig{},
		DatasourceConfig: datasourceConfig,
	}
	var discoverer discovery.Discoverer
	datatsourceStore, _ := store.NewDatasourceStore(datasourceConfig, discoverer)
	remoteStore := store.NewRemoteStore(&http.Client{}, datasourceConfig.QueryTimeout)
	stores = &store.Stores{
		AlertRuleStore:  &store.AlertRuleStore{},
		DashboardStore:  &store.DashboardStore{},
		DatasourceStore: datatsourceStore,
		UserStore:       &store.UserStore{},
		RemoteStore:     remoteStore,
	}
	handlers = loadHandlers(cfg, stores)
	router = NewRouter(cfg, stores)
}

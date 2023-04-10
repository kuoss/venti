package handler

import (
	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/store"
)

type Handlers struct {
	*alertHandler
	*authHandler
	*configHandler
	*dashboardHandler
	*datasourceHandler
	*remoteHandler
}

func loadHandlers(cfg *model.Config, stores *store.Stores) *Handlers {
	return &Handlers{
		NewAlertHandler(stores.AlertRuleStore),
		NewAuthHandler(stores.UserStore),
		NewConfigHandler(cfg),
		NewDashboardHandler(stores.DashboardStore),
		NewDatasourceHandler(stores.DatasourceStore, stores.RemoteStore),
		NewRemoteHandler(stores.DatasourceStore, stores.RemoteStore),
	}
}

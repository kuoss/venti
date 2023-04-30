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
	*probeHandler
	*remoteHandler
}

func loadHandlers(cfg *model.Config, stores *store.Stores) *Handlers {
	return &Handlers{
		NewAlertHandler(stores.AlertRuleStore),
		NewAuthHandler(stores.UserStore),
		NewConfigHandler(cfg),
		NewDashboardHandler(stores.DashboardStore),
		NewDatasourceHandler(stores.DatasourceStore, stores.RemoteStore),
		NewProbeHandler(),
		NewRemoteHandler(stores.DatasourceStore, stores.RemoteStore),
	}
}

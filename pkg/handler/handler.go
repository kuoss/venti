package handler

import (
	"github.com/kuoss/venti/pkg/handler/remote"
	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/store"
)

type Handlers struct {
	alertHandler      *alertHandler
	authHandler       *authHandler
	configHandler     *configHandler
	dashboardHandler  *dashboardHandler
	datasourceHandler *datasourceHandler
	probeHandler      *probeHandler
	remoteHandler     *remote.RemoteHandler
	statusHandler     *statusHandler
}

func loadHandlers(cfg *model.Config, stores *store.Stores) *Handlers {
	return &Handlers{
		NewAlertHandler(stores.AlertRuleStore),
		NewAuthHandler(stores.UserStore),
		NewConfigHandler(cfg),
		NewDashboardHandler(stores.DashboardStore),
		NewDatasourceHandler(stores.DatasourceStore, stores.RemoteStore),
		NewProbeHandler(),
		remote.New(stores.DatasourceStore, stores.RemoteStore),
		NewStatusHandler(stores.StatusStore),
	}
}

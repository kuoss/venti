package handler

import (
	"github.com/kuoss/venti/pkg/handler/remote"
	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/service"
)

type Handlers struct {
	alertHandler      *alertHandler
	authHandler       *authHandler
	dashboardHandler  *dashboardHandler
	datasourceHandler *datasourceHandler
	probeHandler      *probeHandler
	remoteHandler     *remote.RemoteHandler
	statusHandler     *statusHandler
}

func loadHandlers(cfg *model.Config, services *service.Services) *Handlers {
	return &Handlers{
		NewAlertHandler(services.AlertRuleService, services.AlertingService),
		NewAuthHandler(services.UserService),
		NewDashboardHandler(services.DashboardService),
		NewDatasourceHandler(services.DatasourceService, services.RemoteService),
		NewProbeHandler(),
		remote.New(services.DatasourceService, services.RemoteService),
		NewStatusHandler(services.StatusService),
	}
}

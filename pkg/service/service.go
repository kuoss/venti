package service

import (
	"fmt"
	"net/http"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/venti/pkg/config"
	"github.com/kuoss/venti/pkg/service/alerting"
	"github.com/kuoss/venti/pkg/service/alertrule"
	"github.com/kuoss/venti/pkg/service/dashboard"
	"github.com/kuoss/venti/pkg/service/datasource"
	"github.com/kuoss/venti/pkg/service/discovery"
	"github.com/kuoss/venti/pkg/service/discovery/kubernetes"
	"github.com/kuoss/venti/pkg/service/remote"
	"github.com/kuoss/venti/pkg/service/status"
	"github.com/kuoss/venti/pkg/service/user"
)

type Services struct {
	*alertrule.AlertRuleService
	*dashboard.DashboardService
	*datasource.DatasourceService
	*remote.RemoteService
	*status.StatusService
	*user.UserService
	*alerting.AlertingService
}

func NewServices(cfg *config.Config) (*Services, error) {
	// alertrule
	alertRuleService, err := alertrule.New("")
	if err != nil {
		return nil, fmt.Errorf("new alertRuleService err: %w", err)
	}

	// dashboard
	logger.Debugf("new dashboard Service...")
	dashboardService, err := dashboard.New("etc/dashboards")
	if err != nil {
		return nil, fmt.Errorf("new dashboardService err: %w", err)
	}

	// datasource
	var discoverer discovery.Discoverer
	if cfg.DatasourceConfig.Discovery.Enabled {
		discoverer, err = kubernetes.NewK8sService()
		if err != nil {
			return nil, fmt.Errorf("new k8sService err: %w", err)
		}
	}
	datasourceService, err := datasource.New(&cfg.DatasourceConfig, discoverer)
	if err != nil {
		return nil, fmt.Errorf("new datasourceService err: %w", err)
	}

	// remote
	remoteService := remote.New(&http.Client{}, cfg.DatasourceConfig.QueryTimeout)

	// status
	statusService, err := status.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("new statusService err: %w", err)
	}

	// user
	userService, err := user.New("./data/venti.sqlite3", cfg.UserConfig)
	if err != nil {
		return nil, fmt.Errorf("NewUserService err: %w", err)
	}

	// alerting
	alertingService := alerting.New(cfg, alertRuleService.GetAlertRuleFiles(), datasourceService, remoteService)

	return &Services{
		alertRuleService,
		dashboardService,
		datasourceService,
		remoteService,
		statusService,
		userService,
		alertingService,
	}, nil
}

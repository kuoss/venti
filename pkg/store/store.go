package store

import (
	"fmt"
	"net/http"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/store/alerting"
	"github.com/kuoss/venti/pkg/store/alertrule"
	"github.com/kuoss/venti/pkg/store/dashboard"
	"github.com/kuoss/venti/pkg/store/datasource"
	"github.com/kuoss/venti/pkg/store/discovery"
	"github.com/kuoss/venti/pkg/store/discovery/kubernetes"
	"github.com/kuoss/venti/pkg/store/remote"
	"github.com/kuoss/venti/pkg/store/status"
	"github.com/kuoss/venti/pkg/store/user"
)

type Stores struct {
	*alertrule.AlertRuleStore
	*dashboard.DashboardStore
	*datasource.DatasourceStore
	*remote.RemoteStore
	*status.StatusStore
	*user.UserStore
	*alerting.AlertingStore
}

func NewStores(cfg *model.Config) (*Stores, error) {
	// alertrule
	alertRuleStore, err := alertrule.New("")
	if err != nil {
		return nil, fmt.Errorf("alertrule.New err: %w", err)
	}

	// dashboard
	logger.Debugf("new dashboard store...")
	dashboardStore, err := dashboard.New("etc/dashboards")
	if err != nil {
		return nil, fmt.Errorf("NewDashboardStore err: %w", err)
	}

	// datasource
	var discoverer discovery.Discoverer
	if cfg.DatasourceConfig.Discovery.Enabled {
		discoverer, err = kubernetes.NewK8sStore()
		if err != nil {
			return nil, fmt.Errorf("NewK8sStore err: %w", err)
		}
	}
	datasourceStore, err := datasource.New(&cfg.DatasourceConfig, discoverer)
	if err != nil {
		return nil, fmt.Errorf("NewDatasourceStore err: %w", err)
	}

	// remote
	remoteStore := remote.New(&http.Client{}, cfg.DatasourceConfig.QueryTimeout)

	// status
	storeStore := status.New(cfg)

	// user
	userStore, err := user.New("./data/venti.sqlite3", cfg.UserConfig)
	if err != nil {
		return nil, fmt.Errorf("NewUserStore err: %w", err)
	}

	// alerting
	alertingStore := alerting.New("", alertRuleStore.AlertRuleFiles(), datasourceStore)

	return &Stores{
		alertRuleStore,
		dashboardStore,
		datasourceStore,
		remoteStore,
		storeStore,
		userStore,
		alertingStore,
	}, nil
}

package store

import (
	"fmt"
	"net/http"

	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/store/alerting"
	"github.com/kuoss/venti/pkg/store/alertrule"
	"github.com/kuoss/venti/pkg/store/discovery"
	"github.com/kuoss/venti/pkg/store/discovery/kubernetes"
	"github.com/kuoss/venti/pkg/store/remote"
)

type Stores struct {
	*alerting.AlertingStore
	*alertrule.AlertRuleStore
	*DashboardStore
	*DatasourceStore
	*remote.RemoteStore
	*UserStore
}

func NewStores(cfg *model.Config) (*Stores, error) {

	// alerting
	alertingStore := alerting.New("")

	// alertrule
	alertRuleStore, err := alertrule.New("")
	if err != nil {
		return nil, fmt.Errorf("error on New alertRuleStore: %w", err)
	}

	// dashboard
	dashboardStore, err := NewDashboardStore("etc/dashboards/**/*.y*ml")
	if err != nil {
		return nil, fmt.Errorf("load dashboard configuration failed: %w", err)
	}

	// datasource
	var discoverer discovery.Discoverer
	if cfg.DatasourceConfig.Discovery.Enabled {
		discoverer, err = kubernetes.NewK8sStore()
		if err != nil {
			return nil, fmt.Errorf("load discoverer k8sStore failed: %w", err)
		}
	}
	datasourceStore, err := NewDatasourceStore(cfg.DatasourceConfig, discoverer)
	if err != nil {
		return nil, fmt.Errorf("load datasource configuration failed: %w", err)
	}

	// remote
	remoteStore := remote.New(&http.Client{}, cfg.DatasourceConfig.QueryTimeout)

	// user
	userStore, err := NewUserStore("./data/venti.sqlite3", *cfg.UserConfig)
	if err != nil {
		return nil, fmt.Errorf("load user configuration failed: %w", err)
	}

	return &Stores{
		alertingStore,
		alertRuleStore,
		dashboardStore,
		datasourceStore,
		remoteStore,
		userStore,
	}, nil
}

package store

import (
	"fmt"
	"net/http"

	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/store/alertrule"
	"github.com/kuoss/venti/pkg/store/discovery"
	"github.com/kuoss/venti/pkg/store/discovery/kubernetes"
	"github.com/kuoss/venti/pkg/store/remote"
)

type Stores struct {
	*alertrule.AlertRuleStore
	*DashboardStore
	*DatasourceStore
	*UserStore
	*remote.RemoteStore
}

func LoadStores(cfg *model.Config) (*Stores, error) {
	dashboardStore, err := NewDashboardStore("etc/dashboards/**/*.yaml")
	if err != nil {
		return nil, fmt.Errorf("load dashboard configuration failed: %w", err)
	}

	alertRuleStore, err := alertrule.New("etc/alertrules/*.yaml")
	if err != nil {
		return nil, fmt.Errorf("load alertrule configuration failed: %w", err)
	}
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

	userStore, err := NewUserStore("./data/venti.sqlite3", *cfg.UserConfig)
	if err != nil {
		return nil, fmt.Errorf("load user configuration failed: %w", err)
	}
	remoteStore := remote.New(&http.Client{}, cfg.DatasourceConfig.QueryTimeout)

	return &Stores{
		alertRuleStore,
		dashboardStore,
		datasourceStore,
		userStore,
		remoteStore,
	}, nil
}

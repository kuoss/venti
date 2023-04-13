package store

import (
	"fmt"
	"io"
	"net/http"

  "github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/store/discovery"
	"github.com/kuoss/venti/pkg/store/discovery/kubernetes"
	"gopkg.in/yaml.v3"
)

type Stores struct {
	*AlertRuleStore
	*DashboardStore
	*DatasourceStore
	*UserStore
	*RemoteStore
}

func LoadStores(cfg *model.Config) (*Stores, error) {
	dashboardStore, err := NewDashboardStore("etc/dashboards/**/*.yaml")
	if err != nil {
		return nil, fmt.Errorf("load dashboard configuration failed: %w", err)
	}

	alertRuleStore, err := NewAlertRuleStore("etc/alertrules/*.yaml")
	if err != nil {
		return nil, fmt.Errorf("load alertrule configuration failed: %w", err)
	}
	var discoverer discovery.Discoverer
	if cfg.DatasourcesConfig.Discovery.Enabled {
		discoverer, err = kubernetes.NewK8sStore()
		if err != nil {
			return nil, fmt.Errorf("load discoverer k8sStore failed: %w", err)
		}
	}
	datasourceStore, err := NewDatasourceStore(cfg.DatasourcesConfig, discoverer)
	if err != nil {
		return nil, fmt.Errorf("load datasource configuration failed: %w", err)
	}

	userStore, err := NewUserStore("./data/venti.sqlite3", cfg.UserConfig)
	if err != nil {
		return nil, fmt.Errorf("load user configuration failed: %w", err)
	}
	remoteStore := NewRemoteStore(&http.Client{}, cfg.DatasourcesConfig.QueryTimeout)

	return &Stores{
		alertRuleStore,
		dashboardStore,
		datasourceStore,
		userStore,
		remoteStore,
	}, nil
}

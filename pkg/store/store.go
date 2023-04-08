package store

import (
	"fmt"
	"github.com/kuoss/venti/pkg/model"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
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
	datasourceStore, err := NewDatasourceStore(cfg.DatasourcesConfig)
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

func loadYaml(r io.Reader, y interface{}) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("cannot ReadFile: %w", err)
	}
	if err := yaml.Unmarshal(b, y); err != nil {
		return fmt.Errorf("cannot Unmarshal: %w", err)
	}
	return nil
}

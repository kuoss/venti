package pkg

import (
	"fmt"
	"github.com/kuoss/venti/pkg/store/discovery"
	"github.com/kuoss/venti/pkg/store/discovery/kubernetes"

	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/handler"
	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/query"
	"github.com/kuoss/venti/pkg/store"
)

type Stores struct {
	*store.AlertRuleStore
	*store.DashboardStore
	*store.DatasourceStore
	*store.UserStore
}

func LoadStores(cfg *model.Config) (*Stores, error) {
	dashboardStore, err := store.NewDashboardStore("etc/dashboards/**/*.yaml")
	if err != nil {
		return nil, fmt.Errorf("load dashboard configuration failed: %w", err)
	}

	alertStore, err := store.NewAlertRuleStore("etc/alertrules/*.yaml")
	if err != nil {
		return nil, fmt.Errorf("load alertrule configuration failed: %w", err)
	}

	userStore, err := store.NewUserStore("./data/venti.sqlite3", cfg.UserConfig)
	if err != nil {
		return nil, fmt.Errorf("load user configuration failed: %w", err)
	}

	var discoverer discovery.Discoverer
	if cfg.DatasourcesConfig.Discovery.Enabled {
		discoverer, err = kubernetes.NewK8sStore()
		if err != nil {
			return nil, fmt.Errorf("load discoverer k8sStore failed: %w", err)
		}
	}
	datasourceStore, err := store.NewDatasourceStore(cfg.DatasourcesConfig, discoverer)
	if err != nil {
		return nil, fmt.Errorf("load datasource configuration failed: %w", err)
	}

	return &Stores{
		alertStore,
		dashboardStore,
		datasourceStore,
		userStore,
	}, nil
}

// TODO: add to routerGroup routes
// routerGroup.Use(tokenRequired)

func LoadRouter(s *Stores, config *model.Config) *gin.Engine {

	ch := handler.NewConfigHandler(config)
	ah := handler.NewAlertHandler(s.AlertRuleStore)
	dbh := handler.NewDashboardHandler(s.DashboardStore)
	dsh := handler.NewDatasourceHandler(s.DatasourceStore)
	us := handler.NewAuthHandler(s.UserStore)

	mainLethe, _ := s.DatasourceStore.GetMainDatasourceWithType(model.DatasourceTypeLethe)
	letheQuerier := query.NewHttpQuerier(mainLethe, config.DatasourcesConfig.QueryTimeout)
	lh := handler.NewLetheHandler(letheQuerier)

	mainPrometheus, _ := s.DatasourceStore.GetMainDatasourceWithType(model.DatasourceTypePrometheus)
	prometheusQuerier := query.NewHttpQuerier(mainPrometheus, config.DatasourcesConfig.QueryTimeout)
	ph := handler.NewPrometheusHandler(prometheusQuerier)

	router := gin.New()
	// token not required
	user := router.Group("/user")
	{
		user.POST("/login", us.Login)
		user.POST("/logout", us.Logout)
	}

	// routerGroup routes
	api := router.Group("/api")
	{
		configGroup := api.Group("/config")
		{
			configGroup.GET("/version", ch.Version)
		}

		datasourceGroup := api.Group("/datasource")
		{
			datasourceGroup.GET("/targets", dsh.Targets)
		}

		dashboardGroup := api.Group("/dashboards")
		{
			dashboardGroup.GET("/", dbh.Dashboards)
		}

		alertGroup := api.Group("/alerts")
		{
			alertGroup.GET("/", ah.AlertRuleGroupsList)
		}

		letheGroup := api.Group("/lethe")
		{
			letheGroup.GET("/metadata", lh.Metadata)
			letheGroup.GET("/query", lh.Query)
			letheGroup.GET("/query_range", lh.QueryRange)
		}
		prometheusGroup := api.Group("/prometheus")
		{
			prometheusGroup.GET("/metadata", ph.Metadata)
			prometheusGroup.GET("/query", ph.Query)
			prometheusGroup.GET("/query_range", ph.QueryRange)
		}
	}

	router.Use(handleSPA())
	return router
}

package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/store"
)

func SetupRouter(cfg *model.Config, stores *store.Stores) *gin.Engine {
	handlers := loadHandlers(cfg, stores)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// token not required
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", handlers.authHandler.Login)
		authGroup.POST("/logout", handlers.authHandler.Logout)
	}

	// routerGroup routes
	api := router.Group("/api")
	{
		// TODO: api.Use(tokenRequired)
		alertGroup := api.Group("/alerts")
		{
			alertGroup.GET("/", handlers.alertHandler.AlertRuleGroupsList)
		}
		configGroup := api.Group("/config")
		{
			configGroup.GET("/version", handlers.configHandler.Version)
		}
		dashboardGroup := api.Group("/dashboards")
		{
			dashboardGroup.GET("/", handlers.dashboardHandler.Dashboards)
		}
		datasourceGroup := api.Group("/datasource")
		{
			datasourceGroup.GET("/", handlers.datasourceHandler.Datasources)
			datasourceGroup.GET("/targets", handlers.datasourceHandler.Targets)
		}
		remoteGroup := api.Group("/remote")
		{
			remoteGroup.GET("/metadata", handlers.remoteHandler.Metadata)
			remoteGroup.GET("/query", handlers.remoteHandler.Query)
			remoteGroup.GET("/query_range", handlers.remoteHandler.QueryRange)
		}
	}
	router.Use(handleSPA())
	return router
}

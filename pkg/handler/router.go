package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/service"
)

func NewRouter(cfg *model.Config, services *service.Services) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	handlers := loadHandlers(cfg, services)

	api := router.Group("/api/v1")
	// fixme: api.Use(tokenRequired())
	{
		api.GET("/alerts", handlers.alertHandler.AlertRuleFiles)
		api.GET("/config/version", handlers.configHandler.Version)
		api.GET("/dashboards", handlers.dashboardHandler.Dashboards)

		api.GET("/datasources", handlers.datasourceHandler.Datasources)
		api.GET("/datasources/targets", handlers.datasourceHandler.Targets)

		api.GET("/remote/metadata", handlers.remoteHandler.Metadata)
		api.GET("/remote/query", handlers.remoteHandler.Query)
		api.GET("/remote/query_range", handlers.remoteHandler.QueryRange)

		api.GET("/status/buildinfo", handlers.statusHandler.BuildInfo)
	}

	router.POST("/auth/login", handlers.authHandler.Login)
	router.POST("/auth/logout", handlers.authHandler.Logout)

	router.GET("/-/healthy", handlers.probeHandler.Healthy)
	router.GET("/-/ready", handlers.probeHandler.Ready)

	router.Use(handleSPA())

	return router
}

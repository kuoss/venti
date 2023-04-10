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

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", handlers.authHandler.Login)
		authGroup.POST("/logout", handlers.authHandler.Logout)
	}

	api := router.Group("/api")
	// TODO: api.Use(tokenRequired)
	{
		api.GET("/alerts", handlers.alertHandler.AlertRuleGroupsList)
		api.GET("/config/version", handlers.configHandler.Version)
		api.GET("/dashboards", handlers.dashboardHandler.Dashboards)
		api.GET("/datasources", handlers.datasourceHandler.Datasources)
		api.GET("/datasources/targets", handlers.datasourceHandler.Targets)
		api.GET("/remote/metadata", handlers.remoteHandler.Metadata)
		api.GET("/remote/query", handlers.remoteHandler.Query)
		api.GET("/remote/query_range", handlers.remoteHandler.QueryRange)
	}
	router.Use(handleSPA())
	return router
}

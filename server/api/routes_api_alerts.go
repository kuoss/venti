package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/server/configuration"
)

func routesAPIAlerts(api *gin.RouterGroup) {
	api.GET("/alerts", func(c *gin.Context) {
		c.JSON(200, configuration.GetAlertRuleGroups())
	})
}

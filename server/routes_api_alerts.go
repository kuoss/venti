package server

import (
	"github.com/gin-gonic/gin"
)

func routesAPIAlerts(api *gin.RouterGroup) {
	api.GET("/alerts", func(c *gin.Context) {
		c.JSON(200, GetAlertRuleGroups())
	})
}

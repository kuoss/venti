package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/server/configuration"
	"gopkg.in/yaml.v2"
)

func routesAPIConfig(api *gin.RouterGroup) {
	api.GET("/config/version", func(c *gin.Context) {
		c.JSON(200, configuration.GetConfig().Version)
	})

	api.GET("/config/dashboards", func(c *gin.Context) {
		c.JSON(200, configuration.GetConfig().Dashboards)
	})

	api.GET("/config/dashboards/yaml", func(c *gin.Context) {
		bytes, err := yaml.Marshal(configuration.GetConfig().Dashboards)
		if err != nil {
			c.JSON(500, "cannot marshal")
		}
		c.JSON(200, gin.H{
			"yaml": string(bytes),
		})
	})
}

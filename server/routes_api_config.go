package server

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

func routesAPIConfig(api *gin.RouterGroup) {
	api.GET("/config/version", func(c *gin.Context) {
		c.JSON(200, GetConfig().Version)
	})

	api.GET("/config/dashboards", func(c *gin.Context) {
		c.JSON(200, GetConfig().Dashboards)
	})

	api.GET("/config/dashboards/yaml", func(c *gin.Context) {
		bytes, err := yaml.Marshal(GetConfig().Dashboards)
		if err != nil {
			c.JSON(500, "cannot marshal")
		}
		c.JSON(200, gin.H{
			"yaml": string(bytes),
		})
	})
}

package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/server/configuration"
	"net/http"
)

type configHandler struct {
	*configuration.Config
}

func (ch *configHandler) Version(c *gin.Context) {
	c.JSON(http.StatusOK, ch.Config.Version)
	return
}

func (ch *configHandler) Dashboards(c *gin.Context) {

	//dashboards not in config now
	//c.JSON(http.StatusOK,ch.Config.)
	return
}

/*
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

*/

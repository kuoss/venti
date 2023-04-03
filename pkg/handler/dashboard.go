package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/store"
	"net/http"
)

// todo moved from config handler *should* modify web router path
// 1. /config/dashboards -> GET /dashboards
// 2. /config/dashboards/yaml -> should we return yaml bytes?

type dashboardHandler struct {
	*store.DashboardStore
}

func NewDashboardHandler(ds *store.DashboardStore) *dashboardHandler {
	return &dashboardHandler{ds}
}

//GET /dashboards
func (dh *dashboardHandler) Dashboards(c *gin.Context) {
	c.JSON(http.StatusOK, dh.DashboardStore.Dashboards())
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

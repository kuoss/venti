package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/store/dashboard"
)

type dashboardHandler struct {
	dashboardStore *dashboard.DashboardStore
}

func NewDashboardHandler(s *dashboard.DashboardStore) *dashboardHandler {
	return &dashboardHandler{s}
}

// GET /dashboards
func (h *dashboardHandler) Dashboards(c *gin.Context) {
	c.JSON(http.StatusOK, h.dashboardStore.Dashboards())
}

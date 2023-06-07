package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/service/dashboard"
)

type dashboardHandler struct {
	dashboardService *dashboard.DashboardService
}

func NewDashboardHandler(s *dashboard.DashboardService) *dashboardHandler {
	return &dashboardHandler{s}
}

// GET /dashboards
func (h *dashboardHandler) Dashboards(c *gin.Context) {
	c.JSON(http.StatusOK, h.dashboardService.Dashboards())
}

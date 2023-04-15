package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/store"
)

type dashboardHandler struct {
	dashboardStore *store.DashboardStore
}

func NewDashboardHandler(s *store.DashboardStore) *dashboardHandler {
	return &dashboardHandler{s}
}

// GET /dashboards
func (h *dashboardHandler) Dashboards(c *gin.Context) {
	c.JSON(http.StatusOK, h.dashboardStore.Dashboards())
}

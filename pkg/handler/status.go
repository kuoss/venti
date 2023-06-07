package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/service/status"
)

type statusHandler struct {
	statusService *status.StatusService
}

func NewStatusHandler(statusService *status.StatusService) *statusHandler {
	return &statusHandler{statusService}
}

// GET /api/status/buildinfo
func (h *statusHandler) BuildInfo(c *gin.Context) {
	c.JSON(http.StatusOK, h.statusService.BuildInfo())
}

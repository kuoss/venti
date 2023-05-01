package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/store/status"
)

type statusHandler struct {
	statusStore *status.StatusStore
}

func NewStatusHandler(statusStore *status.StatusStore) *statusHandler {
	return &statusHandler{statusStore}
}

// GET /api/status/buildinfo
func (h *statusHandler) BuildInfo(c *gin.Context) {
	c.JSON(http.StatusOK, h.statusStore.BuildInfo())
}

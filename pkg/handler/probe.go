package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type probeHandler struct{}

func NewProbeHandler() *probeHandler {
	return &probeHandler{}
}

func (h *probeHandler) Healthy(c *gin.Context) {
	c.String(http.StatusOK, "Venti is Healthy.\n")
}

func (h *probeHandler) Ready(c *gin.Context) {
	c.String(http.StatusOK, "Venti is Ready.\n")
}

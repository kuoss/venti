package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/service/alerting"
)

type alertHandler struct {
	alertingService *alerting.AlertingService
}

func NewAlertHandler(s *alerting.AlertingService) *alertHandler {
	return &alertHandler{s}
}

func (h *alertHandler) AlertRuleFiles(c *gin.Context) {
	c.JSON(200, h.alertingService.AlertFiles)
}

func (h *alertHandler) SendTestAlert(c *gin.Context) {
	err := h.alertingService.SendTestAlert()
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "error": err.Error()})
	}
	c.JSON(200, gin.H{"status": "success"})
}

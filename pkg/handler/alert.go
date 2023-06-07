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

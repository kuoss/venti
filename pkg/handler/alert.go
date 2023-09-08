package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/service/alerting"
	"github.com/kuoss/venti/pkg/service/alertrule"
)

type alertHandler struct {
	alertRuleService *alertrule.AlertRuleService
	alertingService  *alerting.AlertingService
}

func NewAlertHandler(alertRuleService *alertrule.AlertRuleService, alertingService *alerting.AlertingService) *alertHandler {
	return &alertHandler{alertRuleService, alertingService}
}

func (h *alertHandler) AlertRuleFiles(c *gin.Context) {
	c.JSON(200, h.alertRuleService.GetAlertRuleFiles())
}

func (h *alertHandler) SendTestAlert(c *gin.Context) {
	err := h.alertingService.SendTestAlert()
	if err != nil {
		c.JSON(500, gin.H{"status": "error", "error": err.Error()})
	}
	c.JSON(200, gin.H{"status": "success"})
}

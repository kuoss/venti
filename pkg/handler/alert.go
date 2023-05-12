package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/store/alerting"
)

type alertHandler struct {
	alertingStore *alerting.AlertingStore
}

func NewAlertHandler(s *alerting.AlertingStore) *alertHandler {
	return &alertHandler{s}
}

func (h *alertHandler) AlertRuleFiles(c *gin.Context) {
	c.JSON(200, h.alertingStore.AlertFiles)
}

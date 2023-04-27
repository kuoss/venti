package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/store/alertrule"
)

type alertHandler struct {
	alertRuleStore *alertrule.AlertRuleStore
}

func NewAlertHandler(s *alertrule.AlertRuleStore) *alertHandler {
	return &alertHandler{s}
}

func (h *alertHandler) AlertRuleFiles(c *gin.Context) {
	c.JSON(200, h.alertRuleStore.AlertRuleFiles())
}

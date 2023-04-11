package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/store"
)

type alertHandler struct {
	alertRuleStore *store.AlertRuleStore
}

func NewAlertHandler(s *store.AlertRuleStore) *alertHandler {
	return &alertHandler{s}
}

func (h *alertHandler) AlertRuleFiles(c *gin.Context) {
	c.JSON(200, h.alertRuleStore.AlertRuleFiles())
}

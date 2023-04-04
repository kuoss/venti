package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/store"
)

type alertHandler struct {
	*store.AlertRuleStore
}

func NewAlertHandler(r *store.AlertRuleStore) *alertHandler {
	return &alertHandler{r}
}

func (ah *alertHandler) AlertRuleGroupsList(c *gin.Context) {
	c.JSON(200, ah.AlertRuleStore.RuleGroupsList())
}

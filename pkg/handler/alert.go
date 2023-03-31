package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/store"
)

type alertHandler struct {
	*store.AlertRuleStore
}

func (ah *alertHandler) AlertRuleGroups(c *gin.Context) {
	c.JSON(200, ah.AlertRuleStore.Groups())
}

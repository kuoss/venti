package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/server/configuration"
)

type alertHandler struct {
}

func (ah *alertHandler) AlertRuleGroups(c *gin.Context) {

	//todo get alert groups

	c.JSON(200, configuration.GetAlertRuleGroups())
}

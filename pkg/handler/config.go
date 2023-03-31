package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/configuration"
	"net/http"
)

type configHandler struct {
	*configuration.Config
}

func (ch *configHandler) Version(c *gin.Context) {
	c.JSON(http.StatusOK, ch.Config.Version)
	return
}

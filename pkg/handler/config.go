package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/model"
	"net/http"
)

type configHandler struct {
	*model.Config
}

func NewConfigHandler(config *model.Config) *configHandler {
	return &configHandler{config}
}

func (ch *configHandler) Version(c *gin.Context) {
	c.JSON(http.StatusOK, ch.Config.Version)
}

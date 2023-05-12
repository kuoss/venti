package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/model"
)

type configHandler struct {
	*model.Config
}

func NewConfigHandler(config *model.Config) *configHandler {
	return &configHandler{config}
}

func (h *configHandler) Version(c *gin.Context) {
	c.JSON(http.StatusOK, h.Config.Version)
}

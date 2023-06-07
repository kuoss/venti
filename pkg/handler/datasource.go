package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/handler/api"
	dsService "github.com/kuoss/venti/pkg/service/datasource"
	"github.com/kuoss/venti/pkg/service/remote"
)

type datasourceHandler struct {
	datasourceService *dsService.DatasourceService
	remoteService     *remote.RemoteService
}

func NewDatasourceHandler(datasourceService *dsService.DatasourceService, remoteService *remote.RemoteService) *datasourceHandler {
	return &datasourceHandler{datasourceService, remoteService}
}

// GET /datasources
func (h *datasourceHandler) Datasources(c *gin.Context) {
	c.JSON(http.StatusOK, h.datasourceService.GetDatasources())
}

// GET /datasources/targets
func (h *datasourceHandler) Targets(c *gin.Context) {
	var results []string
	for _, datasource := range h.datasourceService.GetDatasources() {
		_, body, err := h.remoteService.GET(c.Request.Context(), &datasource, remote.ActionTargets, "state=active")
		if err != nil {
			body = fmt.Sprintf(`{"status":"error","errorType":"%s","error":%q}`, api.ErrorExec, err.Error())
		}
		results = append(results, body)
	}
	c.JSON(http.StatusOK, results)
}

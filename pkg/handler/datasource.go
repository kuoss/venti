package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/handler/api"
	"github.com/kuoss/venti/pkg/store"
	"github.com/kuoss/venti/pkg/store/remote"
)

type datasourceHandler struct {
	datasourceStore *store.DatasourceStore
	remoteStore     *remote.RemoteStore
}

func NewDatasourceHandler(datasourceStore *store.DatasourceStore, remoteStore *remote.RemoteStore) *datasourceHandler {
	return &datasourceHandler{datasourceStore, remoteStore}
}

// GET /datasources
func (h *datasourceHandler) Datasources(c *gin.Context) {
	c.JSON(http.StatusOK, h.datasourceStore.GetDatasources())
}

// GET /datasources/targets
func (h *datasourceHandler) Targets(c *gin.Context) {
	var results []string
	for _, datasource := range h.datasourceStore.GetDatasources() {
		_, body, err := h.remoteStore.GET(c.Request.Context(), &datasource, remote.ActionTargets, "state=active")
		if err != nil {
			body = fmt.Sprintf(`{"status":"error","errorType":"%s","error":%q}`, api.ErrorExec, err.Error())
		}
		results = append(results, body)
	}
	c.JSON(http.StatusOK, results)
}

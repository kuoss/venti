package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/store"
)

type datasourceHandler struct {
	datasourceStore *store.DatasourceStore
	remoteStore     *store.RemoteStore
}

func NewDatasourceHandler(datasourceStore *store.DatasourceStore, remoteStore *store.RemoteStore) *datasourceHandler {
	return &datasourceHandler{datasourceStore, remoteStore}
}

// GET /datasources
func (dh *datasourceHandler) Datasources(c *gin.Context) {
	c.JSON(http.StatusOK, dh.datasourceStore.GetDatasources())
}

// GET /datasources/targets
func (dh *datasourceHandler) Targets(c *gin.Context) {
	var results []string
	for _, datasource := range dh.datasourceStore.GetDatasources() {
		result, err := dh.remoteStore.Get(c.Request.Context(), datasource, "targets", "state=active")
		if err != nil {
			result = fmt.Sprintf(`{"status":"error","errorType":"%s"}`, err.Error())
		}
		results = append(results, result)
	}
	c.JSON(http.StatusOK, results)
}

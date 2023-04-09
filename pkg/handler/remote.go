package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/store"
	"net/url"
	"strconv"
)

type remoteHandler struct {
	datasourceStore *store.DatasourceStore
	remoteStore     *store.RemoteStore
}

func NewRemoteHandler(datasourceStore *store.DatasourceStore, remoteStore *store.RemoteStore) *remoteHandler {
	return &remoteHandler{
		datasourceStore,
		remoteStore,
	}
}

// GET /api/remote/metadata
func (h *remoteHandler) Metadata(c *gin.Context) {
	h.remoteAction(c, "metadata", "")
}

// GET /api/remote/query
func (h *remoteHandler) Query(c *gin.Context) {
	values := url.Values{}
	values.Set("time", c.Query("time"))
	values.Set("timeout", c.Query("timeout"))
	values.Set("query", c.Query("query"))
	h.remoteAction(c, "query", values.Encode())
}

// GET /api/remote/query_range
func (h *remoteHandler) QueryRange(c *gin.Context) {
	values := url.Values{}
	values.Set("start", c.Query("start"))
	values.Set("end", c.Query("end"))
	values.Set("step", c.Query("step"))
	values.Set("timeout", c.Query("timeout"))
	values.Set("query", c.Query("query"))
	h.remoteAction(c, "query_range", values.Encode())
}

func (h *remoteHandler) remoteAction(c *gin.Context, action string, rawQuery string) {
	datasource, err := h.getDatasourceWithParams(c.Query("dsid"), c.Query("dstype"))
	if err != nil {
		responseAPIError(c, &apiError{errorInternal, fmt.Errorf("error on getDatasourceWithParams: %w", err)})
		return
	}
	result, err := h.remoteStore.Get(c.Request.Context(), datasource, action, rawQuery)
	if err != nil {
		responseAPIError(c, &apiError{errorInternal, fmt.Errorf("error on remoteStore.Get: %w", err)})
		return
	}
	c.String(200, result)
}

// Select and return the datasource corresponding to the dsID or dsType parameter
func (h *remoteHandler) getDatasourceWithParams(dsID string, dsType string) (model.Datasource, error) {
	if dsID == "" && dsType == "" {
		return model.Datasource{}, errors.New("either dsID or dsType must be specified")
	}
	// If there is a dsID, return the datasource of the corresponding index
	if dsID != "" {
		dsIndex, err := strconv.Atoi(dsID)
		if err != nil {
			return model.Datasource{}, fmt.Errorf("invalid dsid: %w", err)
		}
		datasource, err := h.datasourceStore.GetDatasourceByIndex(dsIndex)
		if err != nil {
			return model.Datasource{}, fmt.Errorf("error on GetDatasourceByIndex: %w", err)
		}
		return datasource, nil
	}
	// The following handles cases where there is no dsID...
	// Invalid if dsType is neither lethe nor prometheus
	if dsType != string(model.DatasourceTypeLethe) && dsType != string(model.DatasourceTypePrometheus) {
		return model.Datasource{}, errors.New("invalid dstype")
	}
	// Returns the main datasource for the requested dsType
	datasource, err := h.datasourceStore.GetMainDatasourceByType(model.DatasourceType(dsType))
	if err != nil {
		return model.Datasource{}, fmt.Errorf("error on GetMainDatasourceByType: %w", err)
	}
	return datasource, nil
}

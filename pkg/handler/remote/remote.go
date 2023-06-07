package remote

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/handler/api"
	"github.com/kuoss/venti/pkg/model"
	dsService "github.com/kuoss/venti/pkg/service/datasource"
	"github.com/kuoss/venti/pkg/service/remote"
)

type RemoteHandler struct {
	datasourceService *dsService.DatasourceService
	remoteService     *remote.RemoteService
}

func New(datasourceService *dsService.DatasourceService, remoteService *remote.RemoteService) *RemoteHandler {
	return &RemoteHandler{
		datasourceService,
		remoteService,
	}
}

// GET /api/remote/metadata
func (h *RemoteHandler) Metadata(c *gin.Context) {
	h.remoteAction(c, remote.ActionMetadata, "")
}

// GET /api/remote/query
func (h *RemoteHandler) Query(c *gin.Context) {
	values := url.Values{}
	values.Set("time", c.Query("time"))
	values.Set("timeout", c.Query("timeout"))
	values.Set("query", c.Query("query"))
	values.Set("logFormat", c.Query("logFormat"))
	h.remoteAction(c, remote.ActionQuery, values.Encode())
}

// GET /api/remote/query_range
func (h *RemoteHandler) QueryRange(c *gin.Context) {
	values := url.Values{}
	values.Set("start", c.Query("start"))
	values.Set("end", c.Query("end"))
	values.Set("step", c.Query("step"))
	values.Set("timeout", c.Query("timeout"))
	values.Set("query", c.Query("query"))
	values.Set("logFormat", c.Query("logFormat"))
	h.remoteAction(c, remote.ActionQueryRange, values.Encode())
}

func (h *RemoteHandler) remoteAction(c *gin.Context, action remote.Action, rawQuery string) {
	datasource, err := h.getDatasourceWithParams(c.Query("dsid"), c.Query("dstype"))
	if err != nil {
		api.ResponseError(c, api.ErrorInternal, fmt.Errorf("error on getDatasourceWithParams: %w", err))
		return
	}
	code, body, err := h.remoteService.GET(c.Request.Context(), &datasource, action, rawQuery)
	if err != nil {
		api.ResponseError(c, api.ErrorInternal, fmt.Errorf("error on GET: %w", err))
		return
	}
	c.String(code, body)
}

// Select and return the datasource corresponding to the dsID or dsType parameter
func (h *RemoteHandler) getDatasourceWithParams(dsID string, dsType string) (model.Datasource, error) {
	if dsID == "" && dsType == "" {
		return model.Datasource{}, errors.New("either dsID or dsType must be specified")
	}

	// If there is a dsID, return the datasource of the corresponding index
	if dsID != "" {
		dsIndex, err := strconv.Atoi(dsID)
		if err != nil {
			return model.Datasource{}, fmt.Errorf("invalid dsid: %w", err)
		}
		datasource, err := h.datasourceService.GetDatasourceByIndex(dsIndex)
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
	datasource, err := h.datasourceService.GetMainDatasourceByType(model.DatasourceType(dsType))
	if err != nil {
		return model.Datasource{}, fmt.Errorf("error on GetMainDatasourceByType: %w", err)
	}
	return datasource, nil
}

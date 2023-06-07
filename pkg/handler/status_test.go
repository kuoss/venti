package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/service/status"
	"github.com/stretchr/testify/assert"
)

var (
	statusHandler1 *statusHandler
	statusRouter   *gin.Engine
)

func init() {
	statusService := status.New(&model.Config{
		Version:          "test",
		DatasourceConfig: model.DatasourceConfig{},
		UserConfig:       model.UserConfig{},
	})

	statusHandler1 = NewStatusHandler(statusService)
	statusRouter = gin.New()
	statusRouter.GET("/api/v1/status/buildinfo", statusHandler1.BuildInfo)
}

func TestNewstatusHandler(t *testing.T) {
	assert.NotEmpty(t, statusHandler1)
}

func TestBuildInfo(t *testing.T) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/v1/status/buildinfo", nil)
	assert.NoError(t, err)
	statusRouter.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	// {"version":"test","goVersion":"go1.20.3"}
	assert.Regexp(t, `{"version":"test","goVersion":"go[0-9.]+"}`, w.Body.String())
}

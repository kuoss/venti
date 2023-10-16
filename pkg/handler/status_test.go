package handler

import (
	"net/http"
	"net/http/httptest"
	"runtime"
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
	statusService, err := status.New(&model.Config{
		AppInfo:          model.AppInfo{Version: "test"},
		DatasourceConfig: model.DatasourceConfig{},
		UserConfig:       model.UserConfig{},
	})
	if err != nil {
		panic(err)
	}

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
	goVersion := runtime.Version()
	assert.Equal(t, `{"data":{"version":"test","revision":"(TBD)","branch":"(TBD)","buildUser":"(TBD)","buildDate":"(TBD)","goVersion":"`+goVersion+`"},"status":"success"}`, w.Body.String())
}

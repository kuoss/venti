package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	probeHandler1 *probeHandler
	probeRouter   *gin.Engine
)

func init() {
	probeHandler1 = NewProbeHandler()
	probeRouter = gin.New()
	probeRouter.GET("/-/healthy", probeHandler1.Healthy)
	probeRouter.GET("/-/ready", probeHandler1.Ready)
}

func TestNewProbeHandler(t *testing.T) {
	assert.NotNil(t, probeHandler1)
}

func TestHealthy(t *testing.T) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/-/healthy", nil)
	assert.NoError(t, err)
	probeRouter.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Venti is Healthy.\n", w.Body.String())
}

func TestReady(t *testing.T) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/-/ready", nil)
	assert.NoError(t, err)
	probeRouter.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Venti is Ready.\n", w.Body.String())
}

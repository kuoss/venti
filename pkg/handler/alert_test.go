package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestNewAlertHandler(t *testing.T) {
	require.NotEmpty(t, handlers.alertHandler)
}

func TestAlertRuleFiles(t *testing.T) {
	alertHandler1 := handlers.alertHandler
	r := gin.Default()
	r.GET("/", alertHandler1.Alerts)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, 200, w.Code)
	want := `{"data":[{"datasourceSelector":{"system":"","type":"prometheus"},"groupLabels":{"rulefile":"sample-v3","severity":"silence"},"alertingRules":[{"rule":{"alert":"S00-AlwaysOn","expr":"vector(1234)","for":"0s","labels":{"hello":"world"},"annotations":{"summary":"AlwaysOn value={{ $value }}"}}},{"rule":{"alert":"S01-Monday","expr":"day_of_week() == 1 and hour() \u003c 2","for":"0s","annotations":{"summary":"Monday"}}},{"rule":{"alert":"S02-NewNamespace","expr":"time() - kube_namespace_created \u003c 120","for":"0s","annotations":{"summary":"labels={{ $labels }} namespace={{ $labels.namespace }} value={{ $value }}"}}},{"rule":{"alert":"PodNotHealthy","expr":"sum by (namespace, pod) (kube_pod_status_phase{phase=~\"Pending|Unknown|Failed\"}) \u003e 0","for":"3s","annotations":{"summary":"{{ $labels.namespace }}/{{ $labels.pod }}"}}}]}],"status":"success"}`
	require.Equal(t, want, w.Body.String())
}

func TestSendTestAlert(t *testing.T) {
	alertHandler1 := handlers.alertHandler
	r := gin.Default()
	r.GET("/", alertHandler1.SendTestAlert)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, 200, w.Code)
	want := `{"status":"success"}`
	require.Equal(t, want, w.Body.String())
}

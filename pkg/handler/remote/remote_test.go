package remote

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kuoss/common/logger"
	ms "github.com/kuoss/venti/pkg/mock/servers"
	"github.com/kuoss/venti/pkg/model"
	dsService "github.com/kuoss/venti/pkg/service/datasource"
	"github.com/kuoss/venti/pkg/service/discovery"
	"github.com/kuoss/venti/pkg/service/remote"
	"github.com/stretchr/testify/assert"
)

var (
	servers        *ms.Servers
	remoteHandler1 *RemoteHandler
	remoteRouter   *gin.Engine
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func shutdown() {
	servers.Close()
}

func setup() {
	logger.SetLevel(logger.DebugLevel)
	servers = ms.New(ms.Requirements{
		{Type: ms.TypePrometheus, Name: "prometheus1", IsMain: true},
		{Type: ms.TypePrometheus, Name: "prometheus2", IsMain: false},
		{Type: ms.TypeLethe, Name: "lethe1", IsMain: true},
		{Type: ms.TypeLethe, Name: "lethe2", IsMain: false},
	})
	var discoverer discovery.Discoverer
	datasourceService, err := dsService.New(&model.DatasourceConfig{Datasources: servers.GetDatasources()}, discoverer)
	if err != nil {
		panic(err)
	}
	remoteService := remote.New(&http.Client{}, 30*time.Second)
	remoteHandler1 = New(datasourceService, remoteService)

	// router
	remoteRouter = gin.New()
	api := remoteRouter.Group("/api")
	api.GET("/remote/metadata", remoteHandler1.Metadata)
	api.GET("/remote/query", remoteHandler1.Query)
	api.GET("/remote/query_range", remoteHandler1.QueryRange)
}

func TestNew(t *testing.T) {
	assert.NotEmpty(t, remoteHandler1.datasourceService)

	assert.NotEmpty(t, remoteHandler1.remoteService)

	var ds model.Datasource
	var err error

	ds, err = remoteHandler1.datasourceService.GetDatasourceByIndex(0)
	assert.NoError(t, err)
	assert.Equal(t, "prometheus1", ds.Name)

	ds, err = remoteHandler1.datasourceService.GetDatasourceByIndex(3)
	assert.NoError(t, err)
	assert.Equal(t, "lethe2", ds.Name)
}

func TestBuildInfo(t *testing.T) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/remote/metadata?dsName=prometheus1", nil)
	assert.NoError(t, err)

	remoteRouter.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"status":"success","data":{"apiserver_audit_event_total":[{"type":"counter","help":"[ALPHA] Counter of audit events generated and sent to the audit backend.","unit":""}]}}`, w.Body.String())
}

func TestMetadata(t *testing.T) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/remote/metadata?dsName=prometheus1", nil)
	assert.NoError(t, err)

	remoteRouter.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"status":"success","data":{"apiserver_audit_event_total":[{"type":"counter","help":"[ALPHA] Counter of audit events generated and sent to the audit backend.","unit":""}]}}`, w.Body.String())
}

func TestQuery(t *testing.T) {
	testCases := []struct {
		rawQuery string
		wantCode int
		wantBody string
	}{
		{
			"",
			405, `{"status":"error","errorType":"bad_data","error":"invalid parameter \"query\": 1:1: parse error: no expression found in input"}`,
		},
		{
			"query=up",
			200, `{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"up","job":"prometheus","instance":"localhost:9090"},"value":[1435781451.781,"1"]},{"metric":{"__name__":"up","job":"prometheus2","instance2":"localhost:9092"},"value":[1435781451.781,"1"]}]}}`,
		},
		{
			"query=not_exists",
			200, `{"status":"success","data":{"resultType":"vector","result":[]}}`,
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/api/remote/query?dsName=prometheus1&"+tc.rawQuery, nil)
			assert.NoError(t, err)

			remoteRouter.ServeHTTP(w, req)
			assert.Equal(t, tc.wantCode, w.Code)
			assert.Equal(t, tc.wantBody, w.Body.String())
		})
	}
}

func TestQueryRange(t *testing.T) {
	testCases := []struct {
		rawQuery string
		wantCode int
		wantBody string
	}{
		{
			"",
			405, `{"status":"error","errorType":"bad_data","error":"invalid parameter \"start\": cannot parse \"\" to a valid timestamp"}`,
		},
		{
			"query=up",
			405, `{"status":"error","errorType":"bad_data","error":"invalid parameter \"start\": cannot parse \"\" to a valid timestamp"}`,
		},
		{
			"query=not_exists",
			405, `{"status":"error","errorType":"bad_data","error":"invalid parameter \"start\": cannot parse \"\" to a valid timestamp"}`,
		},
		{
			"query=not_exists&start=2015-07-01T20:10:30.781Z",
			405, `{"status":"error","errorType":"bad_data","error":"invalid parameter \"end\": cannot parse \"\" to a valid timestamp"}`,
		},
		{
			"query=not_exists&start=2015-07-01T20:10:30.781Z&end=2015-07-01T20:11:00.781Z",
			405, `{"status":"error","errorType":"bad_data","error":"invalid parameter \"step\": cannot parse \"\" to a valid duration"}`,
		},
		{
			"query=not_exists&start=2015-07-01T20:10:30.781Z&end=2015-07-01T20:11:00.781Z&step=15s",
			200, `{"status":"success","data":{"resultType":"matrix","result":[]}}`,
		},
		{
			"query=up&start=2015-07-01T20:10:30.781Z&end=2015-07-01T20:11:00.781Z&step=15s",
			200, `{"status":"success","data":{"resultType":"matrix","result":[{"metric":{"__name__":"up","job":"prometheus","instance":"localhost:9090"},"values":[[1435781430.781,"1"],[1435781445.781,"1"],[1435781460.781,"1"]]}]}}`,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d - %s", i, tc.rawQuery), func(tt *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/api/remote/query_range?dsName=prometheus1&"+tc.rawQuery, nil)
			assert.NoError(tt, err)

			remoteRouter.ServeHTTP(w, req)
			assert.Equal(tt, tc.wantCode, w.Code)
			assert.Equal(tt, tc.wantBody, w.Body.String())
		})
	}
}

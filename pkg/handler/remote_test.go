package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	ms "github.com/kuoss/venti/pkg/mock/servers"
	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/store"
	"github.com/kuoss/venti/pkg/store/discovery"
	"github.com/kuoss/venti/pkg/store/remote"
	"github.com/stretchr/testify/assert"
)

var (
	servers        *ms.Servers
	remoteHandler1 *remoteHandler
	remoteRouter   *gin.Engine
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	servers = ms.New(ms.Requirements{
		{Type: ms.TypePrometheus, Port: 0, Name: "prometheus1", IsMain: true},
		{Type: ms.TypePrometheus, Port: 0, Name: "prometheus2", IsMain: false},
		{Type: ms.TypeLethe, Port: 0, Name: "lethe1", IsMain: true},
		{Type: ms.TypeLethe, Port: 0, Name: "lethe2", IsMain: false},
	})
	datasourceConfig := model.DatasourceConfig{Datasources: servers.GetDatasources()}
	var discoverer discovery.Discoverer
	datasourceStore, _ := store.NewDatasourceStore(&datasourceConfig, discoverer)
	remoteStore := remote.New(&http.Client{}, 30*time.Second)
	remoteHandler1 = NewRemoteHandler(datasourceStore, remoteStore)

	// router
	remoteRouter = gin.New()
	api := remoteRouter.Group("/api")
	api.GET("/remote/metadata", remoteHandler1.Metadata)
	api.GET("/remote/query", remoteHandler1.Query)
	api.GET("/remote/query_range", remoteHandler1.QueryRange)
}

func shutdown() {
	servers.Close()
}

func TestNewRemoteHandler(t *testing.T) {
	assert.NotNil(t, remoteHandler1.datasourceStore)
	assert.NotZero(t, remoteHandler1.datasourceStore)

	assert.NotNil(t, remoteHandler1.remoteStore)
	assert.NotZero(t, remoteHandler1.remoteStore)

	var ds model.Datasource
	var err error

	ds, err = remoteHandler1.datasourceStore.GetDatasourceByIndex(0)
	assert.NoError(t, err)
	assert.Equal(t, "prometheus1", ds.Name)

	ds, err = remoteHandler1.datasourceStore.GetDatasourceByIndex(3)
	assert.NoError(t, err)
	assert.Equal(t, "lethe2", ds.Name)
}

func TestMetadata(t *testing.T) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/remote/metadata?dsid=0", nil)
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
			200, `{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"up","job":"prometheus","instance":"localhost:9090"},"value":[1435781451.781,"1"]}]}}`,
		},
		{
			"query=not_exists",
			200, `{"status":"success","data":{"resultType":"vector","result":[]}}`,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d - %s", i, tc.rawQuery), func(tt *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/api/remote/query?dsid=0&"+tc.rawQuery, nil)
			assert.NoError(tt, err)

			remoteRouter.ServeHTTP(w, req)
			assert.Equal(tt, tc.wantCode, w.Code)
			assert.Equal(tt, tc.wantBody, w.Body.String())
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
			req, err := http.NewRequest("GET", "/api/remote/query_range?dsid=0&"+tc.rawQuery, nil)
			assert.NoError(tt, err)

			remoteRouter.ServeHTTP(w, req)
			assert.Equal(tt, tc.wantCode, w.Code)
			assert.Equal(tt, tc.wantBody, w.Body.String())
		})
	}
}

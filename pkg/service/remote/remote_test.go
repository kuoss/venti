package remote

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/kuoss/common/logger"
	ms "github.com/kuoss/venti/pkg/mock/servers"
	"github.com/kuoss/venti/pkg/model"
	"github.com/stretchr/testify/assert"
)

var (
	servers       *ms.Servers
	datasources   []model.Datasource
	remoteService *RemoteService
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	logger.SetLevel(logger.DebugLevel)
	err := os.Chdir("../..")
	if err != nil {
		panic(err)
	}
	servers = ms.New(ms.Requirements{
		{Type: ms.TypeAlertmanager, Name: "alertmanager1"},
		{Type: ms.TypeLethe, Name: "lethe1", IsMain: true},
		{Type: ms.TypeLethe, Name: "lethe2"},
		{Type: ms.TypePrometheus, Name: "prometheus1", IsMain: true},
		{Type: ms.TypePrometheus, Name: "prometheus2", BasicAuth: true},
		{Type: ms.TypePrometheus, Name: "prometheus3"},
	})
	datasources = servers.GetDatasources()
	remoteService = New(&http.Client{}, 15*time.Second)
}

func shutdown() {
	servers.Close()
}

func TestNew(t *testing.T) {
	assert.Equal(t, 6, len(servers.Svrs))
	assert.Equal(t, "alertmanager1", servers.Svrs[0].Name)
	assert.Equal(t, "prometheus3", servers.Svrs[5].Name)

	assert.Equal(t, 5, len(datasources))
	assert.Equal(t, "lethe1", datasources[0].Name)
	assert.Equal(t, "prometheus3", datasources[4].Name)

	assert.Equal(t, &http.Client{}, remoteService.httpClient)
	assert.Equal(t, 15*time.Second, remoteService.timeout)
}

func TestGET(t *testing.T) {
	testCases := []struct {
		url       string
		action    Action
		rawQuery  string
		wantCode  int
		wantBody  string
		wantError string
	}{
		// no action
		{
			datasources[4].URL, "", "",
			404, "404 page not found\n", "",
		},
		{
			"http://{@example.com", "", "",
			0, "", `error on Parse: parse "http://{@example.com": net/url: invalid userinfo`,
		},
		{
			"wrongURL", "", "",
			0, "", `error on Do: Get "": unsupported protocol scheme ""`,
		},
		{
			"http://0.0.0.0:1111", "", "",
			0, "", `error on Do: Get "http://0.0.0.0:1111": dial tcp 0.0.0.0:1111: connect: connection refused`,
		},
		{
			"http://127.0.0.1:99999", "", "",
			0, "", `error on Do: Get "http://127.0.0.1:99999": dial tcp: address 99999: invalid port`,
		},
		// metadata
		{
			datasources[4].URL, ActionMetadata, "",
			200, `{"status":"success","data":{"apiserver_audit_event_total":[{"type":"counter","help":"[ALPHA] Counter of audit events generated and sent to the audit backend.","unit":""}]}}`,
			"",
		},
		// query
		{
			datasources[4].URL, ActionQuery, "",
			405, `{"status":"error","errorType":"bad_data","error":"invalid parameter \"query\": 1:1: parse error: no expression found in input"}`,
			"",
		},
		{
			datasources[4].URL, ActionQuery, "query=up",
			200, `{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"up","job":"prometheus","instance":"localhost:9090"},"value":[1435781451.781,"1"]}]}}`,
			"",
		},
		{
			datasources[4].URL, ActionQuery, "query=not_exists",
			200, `{"status":"success","data":{"resultType":"vector","result":[]}}`,
			"",
		},
		{
			servers.GetServersByType(ms.TypeAlertmanager)[0].URL, ActionQuery, "query=up",
			404, "404 page not found\n",
			"",
		},
		// query_range
		{
			datasources[4].URL, ActionQueryRange, "",
			405, `{"status":"error","errorType":"bad_data","error":"invalid parameter \"start\": cannot parse \"\" to a valid timestamp"}`,
			"",
		},
		{
			datasources[4].URL, ActionQueryRange, "query=up",
			405, `{"status":"error","errorType":"bad_data","error":"invalid parameter \"start\": cannot parse \"\" to a valid timestamp"}`,
			"",
		},
		{
			datasources[4].URL, ActionQueryRange, "query=not_exists",
			405, `{"status":"error","errorType":"bad_data","error":"invalid parameter \"start\": cannot parse \"\" to a valid timestamp"}`,
			"",
		},
		{
			datasources[4].URL, ActionQueryRange, "query=not_exists&start=2015-07-01T20:10:30.781Z",
			405, `{"status":"error","errorType":"bad_data","error":"invalid parameter \"end\": cannot parse \"\" to a valid timestamp"}`,
			"",
		},
		{
			datasources[4].URL, ActionQueryRange, "query=not_exists&start=2015-07-01T20:10:30.781Z&end=2015-07-01T20:11:00.781Z",
			405, `{"status":"error","errorType":"bad_data","error":"invalid parameter \"step\": cannot parse \"\" to a valid duration"}`,
			"",
		},
		{
			datasources[4].URL, ActionQueryRange, "query=not_exists&start=2015-07-01T20:10:30.781Z&end=2015-07-01T20:11:00.781Z&step=15s",
			200, `{"status":"success","data":{"resultType":"matrix","result":[]}}`,
			"",
		},
		{
			datasources[4].URL, ActionQueryRange, "query=up&start=2015-07-01T20:10:30.781Z&end=2015-07-01T20:11:00.781Z&step=15s",
			200, `{"status":"success","data":{"resultType":"matrix","result":[{"metric":{"__name__":"up","job":"prometheus","instance":"localhost:9090"},"values":[[1435781430.781,"1"],[1435781445.781,"1"],[1435781460.781,"1"]]}]}}`,
			"",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(tt *testing.T) {
			code, body, err := remoteService.GET(context.TODO(), &model.Datasource{URL: tc.url}, tc.action, tc.rawQuery)
			if tc.wantError == "" {
				assert.NoError(tt, err)
			} else {
				assert.EqualError(tt, err, tc.wantError)
			}
			assert.Equal(tt, tc.wantCode, code)
			assert.Equal(tt, tc.wantBody, body)
		})
	}
}

func TestGET_basicAuth(t *testing.T) {
	testCases := []struct {
		datasource *model.Datasource
		wantCode   int
		wantBody   string
		wantError  string
	}{
		{
			&model.Datasource{URL: datasources[3].URL},
			401, "401 unauthorized\n",
			"",
		},
		{
			&model.Datasource{URL: datasources[3].URL, BasicAuth: true, BasicAuthUser: "abc", BasicAuthPassword: "123"},
			200, `{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"up","job":"prometheus","instance":"localhost:9090"},"value":[1435781451.781,"1"]}]}}`,
			"",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(tt *testing.T) {
			code, body, err := remoteService.GET(context.TODO(), tc.datasource, ActionQuery, "query=up")
			if tc.wantError == "" {
				assert.NoError(tt, err)
			} else {
				assert.EqualError(tt, err, tc.wantError)
			}
			assert.Equal(tt, tc.wantCode, code)
			assert.Equal(tt, tc.wantBody, body)
		})
	}
}

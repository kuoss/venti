package remote

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	ms "github.com/kuoss/venti/pkg/mock/servers"
	"github.com/kuoss/venti/pkg/model"
	"github.com/stretchr/testify/assert"
)

var (
	servers     *ms.Servers
	datasources []*model.Datasource
	remoteStore *RemoteStore
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	_ = os.Chdir("../..")
	servers = ms.New(ms.Requirements{
		{Type: ms.TypeAlertmanager, Port: 0, Name: "alertmanager1"},
		{Type: ms.TypeLethe, Port: 0, Name: "lethe1", IsMain: true},
		{Type: ms.TypeLethe, Port: 0, Name: "lethe2"},
		{Type: ms.TypePrometheus, Port: 0, Name: "prometheus1", IsMain: true},
		{Type: ms.TypePrometheus, Port: 0, Name: "prometheus2", BasicAuth: true},
		{Type: ms.TypePrometheus, Port: 0, Name: "prometheus3"},
	})
	datasources = servers.GetDatasources()
	remoteStore = New(&http.Client{}, 15*time.Second)
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

	assert.Equal(t, &http.Client{}, remoteStore.httpClient)
	assert.Equal(t, 15*time.Second, remoteStore.timeout)
}

func TestGET_zero(t *testing.T) {
	datasource := datasources[4]
	action := ""
	code, body, err := remoteStore.GET(context.TODO(), datasource, action, "")
	assert.NoError(t, err)
	assert.Equal(t, 404, code)
	assert.Equal(t, "404 page not found\n", body)
}

func TestGET_error_on_Parse(t *testing.T) {
	datasource := &model.Datasource{}
	action := ""
	var err error

	datasource.URL = "http://{@example.com"
	code, body, err := remoteStore.GET(context.TODO(), datasource, action, "")
	assert.EqualError(t, err, `error on Parse: parse "http://{@example.com": net/url: invalid userinfo`)
	assert.Equal(t, 0, code)
	assert.Equal(t, "", body)
}

func TestGET_error_on_Do(t *testing.T) {
	datasource := &model.Datasource{}
	action := ""

	var code int
	var body string
	var err error

	datasource.URL = "wrongURL"
	code, body, err = remoteStore.GET(context.TODO(), datasource, action, "")
	assert.EqualError(t, err, `error on Do: Get "/api/v1/": unsupported protocol scheme ""`)
	assert.Equal(t, 0, code)
	assert.Equal(t, "", body)

	datasource.URL = "http://0.0.0.0:1111"
	code, body, err = remoteStore.GET(context.TODO(), datasource, action, "")
	assert.EqualError(t, err, `error on Do: Get "http://0.0.0.0:1111/api/v1/": dial tcp 0.0.0.0:1111: connect: connection refused`)
	assert.Equal(t, 0, code)
	assert.Equal(t, "", body)

	datasource.URL = "http://127.0.0.1:99999"
	code, body, err = remoteStore.GET(context.TODO(), datasource, action, "")
	assert.EqualError(t, err, `error on Do: Get "http://127.0.0.1:99999/api/v1/": dial tcp: address 99999: invalid port`)
	assert.Equal(t, 0, code)
	assert.Equal(t, "", body)
}

func TestGET_metadata(t *testing.T) {
	datasource := datasources[4]
	action := "metadata"
	code, body, err := remoteStore.GET(context.TODO(), datasource, action, "")
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, `{"status":"success","data":{"apiserver_audit_event_total":[{"type":"counter","help":"[ALPHA] Counter of audit events generated and sent to the audit backend.","unit":""}]}}`, body)
}

func TestGET_query(t *testing.T) {
	datasource := datasources[4]
	action := "query"
	testCases := []struct {
		rawQuery  string
		wantCode  int
		wantBody  string
		wantError string
	}{
		{
			"",
			405,
			`{"status":"error","errorType":"bad_data","error":"invalid parameter \"query\": 1:1: parse error: no expression found in input"}`,
			"",
		},
		{
			"query=up",
			200,
			`{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"up","job":"prometheus","instance":"localhost:9090"},"value":[1435781451.781,"1"]}]}}`,
			"",
		},
		{
			"query=not_exists",
			200,
			`{"status":"success","data":{"resultType":"vector","result":[]}}`,
			"",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d - %s", i, tc.rawQuery), func(tt *testing.T) {
			code, body, err := remoteStore.GET(context.TODO(), datasource, action, tc.rawQuery)
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
	tempDatasource1 := &model.Datasource{
		URL: servers.Svrs[4].Server.URL,
	}
	tempDatasource2 := &model.Datasource{
		URL:               servers.Svrs[4].Server.URL,
		BasicAuth:         true,
		BasicAuthUser:     "abc",
		BasicAuthPassword: "123",
	}

	testCases := []struct {
		datasource *model.Datasource
		wantCode   int
		wantBody   string
		wantError  string
	}{
		{
			tempDatasource1,
			401,
			"401 unauthorized\n",
			"",
		},
		{
			tempDatasource2,
			200,
			`{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"up","job":"prometheus","instance":"localhost:9090"},"value":[1435781451.781,"1"]}]}}`,
			"",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(tt *testing.T) {
			code, body, err := remoteStore.GET(context.TODO(), tc.datasource, "query", "query=up")
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

func TestGET_query_404(t *testing.T) {
	datasource := &model.Datasource{}
	datasource.URL = servers.GetServersByType(ms.TypeAlertmanager)[0].URL
	action := "query"
	testCases := []struct {
		rawQuery string
		wantCode int
		wantBody string
	}{
		{
			"query=up",
			404,
			"404 page not found\n",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d - %s", i, tc.rawQuery), func(tt *testing.T) {
			code, body, err := remoteStore.GET(context.TODO(), datasource, action, tc.rawQuery)
			assert.Nil(tt, err)
			assert.Equal(tt, tc.wantCode, code)
			assert.Equal(tt, tc.wantBody, body)
		})
	}
}
func TestGET_query_range(t *testing.T) {
	datasource := datasources[4]
	action := "query_range"
	testCases := []struct {
		rawQuery string
		wantCode int
		wantBody string
	}{
		{
			"",
			405,
			`{"status":"error","errorType":"bad_data","error":"invalid parameter \"start\": cannot parse \"\" to a valid timestamp"}`,
		},
		{
			"query=up",
			405,
			`{"status":"error","errorType":"bad_data","error":"invalid parameter \"start\": cannot parse \"\" to a valid timestamp"}`,
		},
		{
			"query=not_exists",
			405,
			`{"status":"error","errorType":"bad_data","error":"invalid parameter \"start\": cannot parse \"\" to a valid timestamp"}`,
		},
		{
			"query=not_exists&start=2015-07-01T20:10:30.781Z",
			405,
			`{"status":"error","errorType":"bad_data","error":"invalid parameter \"end\": cannot parse \"\" to a valid timestamp"}`,
		},
		{
			"query=not_exists&start=2015-07-01T20:10:30.781Z&end=2015-07-01T20:11:00.781Z",
			405,
			`{"status":"error","errorType":"bad_data","error":"invalid parameter \"step\": cannot parse \"\" to a valid duration"}`,
		},
		{
			"query=not_exists&start=2015-07-01T20:10:30.781Z&end=2015-07-01T20:11:00.781Z&step=15s",
			200,
			`{"status":"success","data":{"resultType":"matrix","result":[]}}`,
		},
		{
			"query=up&start=2015-07-01T20:10:30.781Z&end=2015-07-01T20:11:00.781Z&step=15s",
			200,
			`{"status":"success","data":{"resultType":"matrix","result":[{"metric":{"__name__":"up","job":"prometheus","instance":"localhost:9090"},"values":[[1435781430.781,"1"],[1435781445.781,"1"],[1435781460.781,"1"]]}]}}`,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d - %s", i, tc.rawQuery), func(tt *testing.T) {
			code, body, err := remoteStore.GET(context.TODO(), datasource, action, tc.rawQuery)
			assert.NoError(tt, err)
			assert.Equal(tt, tc.wantCode, code)
			assert.Equal(tt, tc.wantBody, body)
		})
	}
}

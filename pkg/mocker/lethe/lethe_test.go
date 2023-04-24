package lethe_test

import (
	"testing"

	"github.com/kuoss/venti/pkg/mocker"
	mockerClient "github.com/kuoss/venti/pkg/mocker/client"
	"github.com/kuoss/venti/pkg/mocker/prometheus"
	"github.com/stretchr/testify/assert"
)

var (
	server *mocker.Server
	client *mockerClient.Client
)

func init() {
	server, _ = prometheus.New(0)
	client = mockerClient.New(server.URL)
}

func Test_api_v1_metadata(t *testing.T) {
	code, body, err := client.GET("/api/v1/metadata", "")
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.JSONEq(t, `{"status":"success","data":{"apiserver_audit_event_total":[{"type":"counter","help":"[ALPHA] Counter of audit events generated and sent to the audit backend.","unit":""}]}}`, body)
}

func Test_api_v1_query(t *testing.T) {
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
	for _, tc := range testCases {
		t.Run(tc.rawQuery, func(tt *testing.T) {
			code, body, err := client.GET("/api/v1/query", tc.rawQuery)
			assert.NoError(tt, err)
			assert.Equal(tt, tc.wantCode, code)
			assert.JSONEq(tt, tc.wantBody, body)
		})
	}
}

func Test_api_v1_query_range(t *testing.T) {
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
			"start=2015-07-01T20:10:30.781Z&end=2015-07-01T20:11:00.781Z&step=15s",
			405, `{"status":"error","errorType":"bad_data","error":"1:1: parse error: no expression found in input"}`,
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
	for _, tc := range testCases {
		t.Run(tc.rawQuery, func(tt *testing.T) {
			code, body, err := client.GET("/api/v1/query_range", tc.rawQuery)
			assert.NoError(t, err)
			assert.Equal(tt, tc.wantCode, code)
			assert.JSONEq(tt, tc.wantBody, body)
		})
	}
}

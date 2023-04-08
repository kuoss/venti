package mock

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mockPrometheus *httptest.Server
)

func init() {
	mockPrometheus = MockPrometheusServer()
}

func TestMetadata(t *testing.T) {
	u, _ := url.Parse(mockPrometheus.URL)
	u.Path = "/api/v1/metadata"
	want := `{"status":"success","data":{"apiserver_audit_event_total":[{"type":"counter","help":"[ALPHA] Counter of audit events generated and sent to the audit backend.","unit":""}]}}`

	res, _ := http.Get(u.String())
	bodyBytes, _ := io.ReadAll(res.Body)
	assert.Equal(t, want, string(bodyBytes))
}

func TestQuery(t *testing.T) {
	u, _ := url.Parse(mockPrometheus.URL)
	u.Path = "/api/v1/query"
	testCases := []struct {
		rawQuery string
		want     string
	}{
		{
			"",
			`{"status":"error","errorType":"bad_data","error":"invalid parameter \"query\": 1:1: parse error: no expression found in input"}`,
		},
		{
			"query=up",
			`{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"up","job":"prometheus","instance":"localhost:9090"},"value":[1435781451.781,"1"]}]}}`,
		},
		{
			"query=not_exists",
			`{"status":"success","data":{"resultType":"vector","result":[]}}`,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d - %s", i, tc.rawQuery), func(tt *testing.T) {
			u.RawQuery = tc.rawQuery
			res, _ := http.Get(u.String())
			bodyBytes, _ := io.ReadAll(res.Body)
			assert.Equal(t, tc.want, string(bodyBytes))
		})
	}
}

func TestQueryRange(t *testing.T) {
	u, _ := url.Parse(mockPrometheus.URL)
	u.Path = "/api/v1/query_range"
	testCases := []struct {
		rawQuery string
		want     string
	}{
		{
			"",
			`{"status":"error","errorType":"bad_data","error":"invalid parameter \"start\": cannot parse \"\" to a valid timestamp"}`,
		},
		{
			"query=up",
			`{"status":"error","errorType":"bad_data","error":"invalid parameter \"start\": cannot parse \"\" to a valid timestamp"}`,
		},
		{
			"query=not_exists",
			`{"status":"error","errorType":"bad_data","error":"invalid parameter \"start\": cannot parse \"\" to a valid timestamp"}`,
		},
		{
			"query=not_exists&start=2015-07-01T20:10:30.781Z",
			`{"status":"error","errorType":"bad_data","error":"invalid parameter \"end\": cannot parse \"\" to a valid timestamp"}`,
		},
		{
			"query=not_exists&start=2015-07-01T20:10:30.781Z&end=2015-07-01T20:11:00.781Z",
			`{"status":"error","errorType":"bad_data","error":"invalid parameter \"step\": cannot parse \"\" to a valid duration"}`,
		},
		{
			"query=not_exists&start=2015-07-01T20:10:30.781Z&end=2015-07-01T20:11:00.781Z&step=15s",
			`{"status":"success","data":{"resultType":"matrix","result":[]}}`,
		},
		{
			"query=up&start=2015-07-01T20:10:30.781Z&end=2015-07-01T20:11:00.781Z&step=15s",
			`{"status":"success","data":{"resultType":"matrix","result":[{"metric":{"__name__":"up","job":"prometheus","instance":"localhost:9090"},"values":[[1435781430.781,"1"],[1435781445.781,"1"],[1435781460.781,"1"]]}]}}`,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d - %s", i, tc.rawQuery), func(tt *testing.T) {
			u.RawQuery = tc.rawQuery
			res, _ := http.Get(u.String())
			bodyBytes, _ := io.ReadAll(res.Body)
			assert.Equal(t, tc.want, string(bodyBytes))
		})
	}
}

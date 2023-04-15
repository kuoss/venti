package store

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kuoss/venti/pkg/mock"
	"github.com/kuoss/venti/pkg/model"
	"github.com/stretchr/testify/assert"
)

var (
	ts          *httptest.Server
	timeout     time.Duration
	remoteStore *RemoteStore
	ctx         context.Context
	datasource  model.Datasource
)

func init() {
	ts = mock.PrometheusServer()
	timeout = 30 * time.Second
	remoteStore = NewRemoteStore(&http.Client{}, timeout)
	ctx = (&http.Request{}).Context()
	datasource = model.Datasource{
		Type: model.DatasourceTypePrometheus,
		URL:  ts.URL,
	}
}

func TestNewRemoteStore(t *testing.T) {
	assert.Equal(t, &http.Client{}, remoteStore.httpClient)
	assert.Equal(t, timeout, remoteStore.timeout)
}

func TestGet_Metadata(t *testing.T) {
	action := "metadata"
	result, err := remoteStore.Get(ctx, datasource, action, "")
	assert.Nil(t, err)
	assert.Equal(t, `{"status":"success","data":{"apiserver_audit_event_total":[{"type":"counter","help":"[ALPHA] Counter of audit events generated and sent to the audit backend.","unit":""}]}}`, result)
}

func TestGet_Query(t *testing.T) {
	action := "query"
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
			result, err := remoteStore.Get(ctx, datasource, action, tc.rawQuery)
			assert.Nil(tt, err)
			assert.Equal(tt, tc.want, result)
		})
	}
}

func TestGet_QueryRange(t *testing.T) {
	action := "query_range"
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
			result, err := remoteStore.Get(ctx, datasource, action, tc.rawQuery)
			assert.Nil(tt, err)
			assert.Equal(tt, tc.want, result)
		})
	}
}

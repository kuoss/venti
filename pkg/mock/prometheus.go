package mock

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
)

func mockResponse(w http.ResponseWriter, code int, body string) {
	w.WriteHeader(code)
	fmt.Fprint(w, body)
}

func MockPrometheusServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/metadata", func(w http.ResponseWriter, r *http.Request) {
		mockResponse(w, 200, `{"status":"success","data":{"apiserver_audit_event_total":[{"type":"counter","help":"[ALPHA] Counter of audit events generated and sent to the audit backend.","unit":""}]}}`)
	})
	mux.HandleFunc("/api/v1/query", func(w http.ResponseWriter, r *http.Request) {
		q, _ := url.ParseQuery(r.URL.RawQuery)
		if q.Get("query") == "" {
			mockResponse(w, 405, `{"status":"error","errorType":"bad_data","error":"invalid parameter \"query\": 1:1: parse error: no expression found in input"}`)
			return
		}
		if q.Get("query") == "up" {
			mockResponse(w, 200, `{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"up","job":"prometheus","instance":"localhost:9090"},"value":[1435781451.781,"1"]}]}}`)
			return
		}
		// metric_not_exists
		mockResponse(w, 200, `{"status":"success","data":{"resultType":"vector","result":[]}}`)
	})
	mux.HandleFunc("/api/v1/query_range", func(w http.ResponseWriter, r *http.Request) {
		q, _ := url.ParseQuery(r.URL.RawQuery)
		if q.Get("start") == "" {
			mockResponse(w, 405, `{"status":"error","errorType":"bad_data","error":"invalid parameter \"start\": cannot parse \"\" to a valid timestamp"}`)
			return
		}
		if q.Get("end") == "" {
			mockResponse(w, 405, `{"status":"error","errorType":"bad_data","error":"invalid parameter \"end\": cannot parse \"\" to a valid timestamp"}`)
			return
		}
		if q.Get("step") == "" {
			mockResponse(w, 405, `{"status":"error","errorType":"bad_data","error":"invalid parameter \"step\": cannot parse \"\" to a valid duration"}`)
			return
		}
		if q.Get("query") == "" {
			mockResponse(w, 405, `{"status":"error","errorType":"bad_data","error":"1:1: parse error: no expression found in input"}`)
			return
		}
		if q.Get("query") == "up" {
			mockResponse(w, 200, `{"status":"success","data":{"resultType":"matrix","result":[{"metric":{"__name__":"up","job":"prometheus","instance":"localhost:9090"},"values":[[1435781430.781,"1"],[1435781445.781,"1"],[1435781460.781,"1"]]}]}}`)
			return
		}
		// metric_not_exists
		mockResponse(w, 200, `{"status":"success","data":{"resultType":"matrix","result":[]}}`)
	})
	return httptest.NewServer(mux)
}

package mock

import (
	"net/http"
	"net/http/httptest"
)

func Alertmanager() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/alerts", func(w http.ResponseWriter, r *http.Request) {
		response(w, 200, `{"status":"success","data":{}}`)
	})
	return httptest.NewServer(mux)
}

package server

import (
	"errors"
	"io"
	"net/http"
)

var client *http.Client

func init() {
	client = &http.Client{Timeout: config.DatasourcesConfig.QueryTimeout}
}

func HTTPGet(url string, params map[string]string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("cannot connect to datasource")
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), err
}

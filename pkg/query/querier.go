package query

import (
	"context"
	"fmt"
	"github.com/kuoss/venti/pkg/configuration"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Querier interface {
	Execute(ctx context.Context, q Query) (string, error)
}

type httpQuerier struct {
	*http.Client
	url        string
	timeout    time.Duration
	datasource configuration.Datasource
}

func NewHttpQuerier(ds configuration.Datasource, timeout time.Duration) *httpQuerier {
	return &httpQuerier{
		Client:     &http.Client{},
		url:        ds.URL,
		timeout:    timeout,
		datasource: ds,
	}
}

func (hq *httpQuerier) Execute(ctx context.Context, q Query) (string, error) {

	// ds := q.Datasource
	// todo: what for?
	/*
		if ds.Type == DatasourceTypeNone {
			var err error
			ds, err = configuration.GetDefaultDatasource(q.DatasourceType)
			if err != nil {
				return "", fmt.Errorf("error on GetDefaultDatasource: %w", err)
			}
		}

	*/

	u, err := url.Parse(hq.url + q.Path)
	if err != nil {
		return "", fmt.Errorf("error on url.Parse: %w", err)
	}

	queryParam := u.Query()
	for key, value := range q.Params {
		queryParam.Set(key, value)
	}

	u.RawQuery = queryParam.Encode()
	ctx, cancel := context.WithTimeout(ctx, hq.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("error on http.NewRequest: %w", err)
	}

	if hq.datasource.BasicAuth {
		req.SetBasicAuth(hq.datasource.BasicAuthUser, hq.datasource.BasicAuthPassword)
	}

	resp, err := hq.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error on client.Do: %w", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error on io.ReadAll: %w", err)
	}
	return string(body), nil
}

// todo string return type?

/*
func (hq *httpQuerier) Get(ctx context.Context, params map[string]string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, hq.url, nil)
	if err != nil {
		return "", err
	}
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	resp, err := hq.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("cannot connect to datasource: %w", err)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read datasource query response: %w", err)
	}
	return string(bodyBytes), err
}
*/

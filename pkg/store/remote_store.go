package store

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/kuoss/venti/pkg/model"
)

type RemoteStore struct {
	httpClient *http.Client
	timeout    time.Duration
}

func NewRemoteStore(httpClient *http.Client, timeout time.Duration) *RemoteStore {
	return &RemoteStore{
		httpClient: httpClient,
		timeout:    timeout,
	}
}

func (r *RemoteStore) Get(ctx context.Context, datasource model.Datasource, action string, rawQuery string) (string, error) {
	u, err := url.Parse(datasource.URL)
	if err != nil {
		return "", fmt.Errorf("error on url.Parse: %w", err)
	}
	u.Path = fmt.Sprintf("/api/v1/%s", action)
	u.RawQuery = rawQuery
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("error on http.NewRequest: %w", err)
	}

	if datasource.BasicAuth {
		req.SetBasicAuth(datasource.BasicAuthUser, datasource.BasicAuthPassword)
	}
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error on client.Do: %w", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error on io.ReadAll: %w", err)
	}
	return string(body), nil
}

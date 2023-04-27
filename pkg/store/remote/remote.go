package remote

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

func New(httpClient *http.Client, timeout time.Duration) *RemoteStore {
	return &RemoteStore{
		httpClient: httpClient,
		timeout:    timeout,
	}
}

func (r *RemoteStore) GET(ctx context.Context, datasource *model.Datasource, action string, rawQuery string) (code int, body string, err error) {
	u, err := url.Parse(datasource.URL)
	if err != nil {
		return 0, "", fmt.Errorf("error on Parse: %w", err)
	}

	u.Path = fmt.Sprintf("/api/v1/%s", action)
	u.RawQuery = rawQuery
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		//go:cover ignore - hardly reachable
		return 0, "", fmt.Errorf("error on NewRequest: %w", err)
	}

	if datasource.BasicAuth {
		// TODO: test cover
		req.SetBasicAuth(datasource.BasicAuthUser, datasource.BasicAuthPassword)
	}
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("error on Do: %w", err)
	}
	code = resp.StatusCode
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		//go:cover ignore - hardly reachable
		return code, "", fmt.Errorf("error on ReadAll: %w", err)
	}
	return code, string(bodyBytes), nil
}

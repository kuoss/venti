package remote

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/venti/pkg/model"
)

type RemoteService struct {
	httpClient *http.Client
	timeout    time.Duration
}

type Action string

const (
	ActionReady      Action = "/-/ready"
	ActionMetadata   Action = "/api/v1/metadata"
	ActionQuery      Action = "/api/v1/query"
	ActionQueryRange Action = "/api/v1/query_range"
	ActionTargets    Action = "/api/v1/targets"
)

func New(httpClient *http.Client, timeout time.Duration) *RemoteService {
	return &RemoteService{
		httpClient: httpClient,
		timeout:    timeout,
	}
}

func (r *RemoteService) GET(ctx context.Context, datasource *model.Datasource, action Action, rawQuery string) (code int, body string, err error) {
	u, err := url.Parse(datasource.URL)
	if err != nil {
		return 0, "", fmt.Errorf("error on Parse: %w", err)
	}

	u.Path = string(action)
	u.RawQuery = rawQuery
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		// unreachable
		return 0, "", fmt.Errorf("NewRequest err: %w", err)
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
		// unreachable
		return code, "", fmt.Errorf("error on ReadAll: %w", err)
	}
	if code != 200 {
		logger.Debugf("code=%d url=%s", code, u.String())
	}
	return code, string(bodyBytes), nil
}

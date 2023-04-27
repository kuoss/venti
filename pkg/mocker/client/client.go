package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	endpoint          string
	basicAuth         bool
	basicAuthUser     string
	basicAuthPassword string
}

func New(endpoint string) *Client {
	return &Client{endpoint: endpoint}
}

func (c *Client) SetBasicAuth(username string, password string) {
	c.basicAuth = true
	c.basicAuthUser = username
	c.basicAuthPassword = password
}

func (c *Client) GET(path string, rawQuery string) (code int, body string, err error) {
	u, err := url.Parse(c.endpoint)
	if err != nil {
		return 0, "", fmt.Errorf("error on Parse: %w", err)
	}

	u.Path = path
	u.RawQuery = rawQuery
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return 0, "", fmt.Errorf("error on NewRequest: %w", err)
	}

	if c.basicAuth {
		req.SetBasicAuth(c.basicAuthUser, c.basicAuthPassword)
	}
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("error on Get: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", fmt.Errorf("error on ReadAll: %w", err)
	}

	return resp.StatusCode, string(bodyBytes), nil
}

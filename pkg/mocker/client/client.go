package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	url string
}

func New(url string) *Client {
	return &Client{url}
}

func (c *Client) GET(path string, rawQuery string) (code int, body string, err error) {
	u, err := url.Parse(c.url)
	if err != nil {
		return 0, "", fmt.Errorf("error on Parse: %w", err)
	}
	u.Path = path
	u.RawQuery = rawQuery
	resp, err := http.Get(u.String())
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

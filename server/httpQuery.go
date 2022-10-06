package server

import (
	"fmt"
)

type HTTPQuery struct {
	Host  string `form:"host,omitempty"`
	Port  int    `form:"port,omitempty"`
	Query string `form:"expr"`
	Time  string `form:"time,omitempty"`
}

type HTTPQueryRange struct {
	Host  string `form:"host,omitempty"`
	Port  int    `form:"port,omitempty"`
	Query string `form:"expr"`
	Start string `form:"start"`
	End   string `form:"end"`
	Step  string `form:"step"`
}

func RunHTTPPrometheusQuery(httpQuery HTTPQuery) (string, error) {
	if httpQuery.Host == "" {
		httpQuery.Host = "prometheus"
	}
	if httpQuery.Port == 0 {
		httpQuery.Port = 9090
	}
	return runHTTPQuery(httpQuery)
}

func RunHTTPPrometheusQueryRange(httpQueryRange HTTPQueryRange) (string, error) {
	if httpQueryRange.Host == "" {
		httpQueryRange.Host = "prometheus"
	}
	if httpQueryRange.Port == 0 {
		httpQueryRange.Port = 9090
	}
	return runHTTPQueryRange(httpQueryRange)
}

func RunHTTPLetheQuery(httpQuery HTTPQuery) (string, error) {
	if httpQuery.Host == "" {
		httpQuery.Host = "lethe"
	}
	if httpQuery.Port == 0 {
		httpQuery.Port = 8080
	}
	return runHTTPQuery(httpQuery)
}

func RunHTTPLetheQueryRange(httpQueryRange HTTPQueryRange) (string, error) {
	if httpQueryRange.Host == "" {
		httpQueryRange.Host = "lethe"
	}
	if httpQueryRange.Port == 0 {
		httpQueryRange.Port = 8080
	}
	return runHTTPQueryRange(httpQueryRange)
}

func runHTTPQuery(httpQuery HTTPQuery) (string, error) {
	// log.Printf("http://%s:%d/api/v1/query query=%s", httpQuery.Host, httpQuery.Port, httpQuery.Query)
	response, err := HTTPGet(fmt.Sprintf("http://%s:%d/api/v1/query", httpQuery.Host, httpQuery.Port), map[string]string{
		"query": httpQuery.Query,
		"time":  httpQuery.Time,
	})
	if err != nil {
		return "", err
	}
	return response, nil
}

func runHTTPQueryRange(httpQueryRange HTTPQueryRange) (string, error) {
	response, err := HTTPGet(fmt.Sprintf("http://%s:%d/api/v1/query_range", httpQueryRange.Host, httpQueryRange.Port), map[string]string{
		"query": httpQueryRange.Query,
		"start": httpQueryRange.Start,
		"end":   httpQueryRange.End,
		"step":  httpQueryRange.Step,
	})
	if err != nil {
		return "", err
	}
	return response, nil
}

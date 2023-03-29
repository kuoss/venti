package server

import (
	"fmt"
	"github.com/kuoss/venti/server/configuration"
	"io"
	"log"
	"net/http"
	"net/url"
)

type PathQuery struct {
	Datasource     Datasource
	DatasourceType DatasourceType
	Path           string
	Params         map[string]string
	Start          string
	End            string
	Step           string
}

type InstantQuery struct {
	Datasource     Datasource     `form:"datasource,omitempty"`
	DatasourceType DatasourceType `form:"datasourceType,omitempty"`
	Expr           string         `form:"expr"`
	Time           string         `form:"time,omitempty"`
}

type RangeQuery struct {
	Datasource     Datasource     `form:"datasource,omitempty"`
	DatasourceType DatasourceType `form:"datasourceType,omitempty"`
	Expr           string         `form:"expr"`
	Start          string         `form:"start,omitempty"`
	End            string         `form:"end,omitempty"`
	Step           string         `form:"step,omitempty"`
}

func (pq PathQuery) execute() (string, error) {
	ds := pq.Datasource
	if ds.Type == DatasourceTypeNone {
		var err error
		ds, err = configuration.GetDefaultDatasource(pq.DatasourceType)
		if err != nil {
			return "", fmt.Errorf("error on GetDefaultDatasource: %w", err)
		}
	}
	u, err := url.Parse(ds.URL + pq.Path)
	if err != nil {
		return "", fmt.Errorf("error on url.Parse: %w", err)
	}
	q := u.Query()
	for key, value := range pq.Params {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()
	log.Println(u.String())
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("error on http.NewRequest: %w", err)
	}
	if ds.BasicAuth {
		req.SetBasicAuth(ds.BasicAuthUser, ds.BasicAuthPassword)
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error on client.Do: %w", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error on io.ReadAll: %w", err)
	}
	return string(body), nil
}

func (iq InstantQuery) execute() (string, error) {
	return PathQuery{
		DatasourceType: iq.DatasourceType,
		Path:           "/handler/v1/query",
		Params: map[string]string{
			"query": iq.Expr,
			"time":  iq.Time,
		},
	}.execute()
}

func (rq RangeQuery) execute() (string, error) {
	return PathQuery{
		DatasourceType: rq.DatasourceType,
		Path:           "/handler/v1/query_range",
		Params: map[string]string{
			"query": rq.Expr,
			"start": rq.Start,
			"end":   rq.End,
			"step":  rq.Step,
		},
	}.execute()
}

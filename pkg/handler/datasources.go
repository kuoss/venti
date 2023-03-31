package handler

import (
	"context"
	"fmt"
	"github.com/kuoss/venti/pkg/configuration"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type datasourceHandler struct {
	*configuration.DatasourcesConfig
}

// datasources
func (dh *datasourceHandler) Datasources(c *gin.Context) {
	// todo check Datasources is pointer type
	c.JSON(http.StatusOK, dh.DatasourcesConfig.Datasources)
}

///datasources/targets

func (dh *datasourceHandler) Targets(c *gin.Context) {

	var bodies []string
	for _, ds := range dh.DatasourcesConfig.Datasources {
		url := fmt.Sprintf("%s/handler/v1/targets?state=active", ds.URL)
		body, err := httpDo(url, *ds)
		if err != nil {
			bodies = append(bodies, fmt.Sprintf(`{"status":"error","errorType":"%s"}`, err.Error()))
		} else {
			b, err := io.ReadAll(body)
			if err != nil {
				log.Printf("%s read body from datasource fail %v", ds.Name, err)
				continue
			}
			bodies = append(bodies, string(b))
		}
	}
	c.JSON(http.StatusOK, bodies)
	return
}

func httpDo(url string, ds configuration.Datasource) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {

	}
	if ds.BasicAuth {
		req.SetBasicAuth(ds.BasicAuthUser, ds.BasicAuthPassword)
	}
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("datasource %s target state response not ok", ds.Name)
	}
	return resp.Body, err
}

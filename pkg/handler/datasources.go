package handler

import (
	"context"
	"fmt"
	"github.com/kuoss/venti/pkg/configuration"
	"github.com/kuoss/venti/pkg/store"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type datasourceHandler struct {
	// configuration.DatasourcesConfig should be used for initializing datasourceStore
	*store.DatasourceStore
}

func NewDatasourceHandler(ds *store.DatasourceStore) *datasourceHandler {
	return &datasourceHandler{ds}
}

// GET datasources
func (dh *datasourceHandler) Datasources(c *gin.Context) {

	// todo check Datasources is pointer type
	c.JSON(http.StatusOK, dh.DatasourceStore.GetDatasources())
}

// GET /datasources/targets

func (dh *datasourceHandler) Targets(c *gin.Context) {

	var bodies []string
	for _, ds := range dh.DatasourceStore.GetDatasources() {
		url := fmt.Sprintf("%s/api/v1/targets?state=active", ds.URL)
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

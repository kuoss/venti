package api

import (
	"fmt"
	"github.com/kuoss/venti/server/configuration"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func routesAPIDatasources(api *gin.RouterGroup) {
	api.GET("/datasources", func(c *gin.Context) {
		datasources := configuration.GetDatasources()
		c.JSON(200, datasources)
	})
	api.GET("/datasources/targets", func(c *gin.Context) {
		datasources := configuration.GetDatasources()
		var bodies []string
		for _, ds := range datasources {
			var err error
			var body string
			apiURL := fmt.Sprintf("%s/api/v1/targets?state=active", ds.URL)
			log.Println(apiURL)
			req, err := http.NewRequest("GET", apiURL, nil)
			if err != nil {
				bodies = append(bodies, `{"status":"error","errorType":"NewRequest"}`)
				continue
			}
			if ds.BasicAuth {
				req.SetBasicAuth(ds.BasicAuthUser, ds.BasicAuthPassword)
			}
			client := http.Client{Timeout: 2 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				body = `{"status":"error","errorType":"timeout"}`
			} else {
				bodyBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					body = `{"status":"error","errorType":"internal"}`
				} else {
					if resp.StatusCode != http.StatusOK {
						body = fmt.Sprintf(`{"status":"error","errorType":"%s"}`, resp.Status)
					} else {
						body = string(bodyBytes)
					}
				}
			}
			bodies = append(bodies, body)
		}
		c.JSON(200, bodies)
	})
}

package server

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func routesAPIDatasources(api *gin.RouterGroup) {
	api.GET("/datasources", func(c *gin.Context) {
		datasources, err := GetDatasources()
		if err != nil {
			c.JSON(500, err)
		}
		c.JSON(200, datasources)
	})
	api.GET("/datasources/targets", func(c *gin.Context) {
		datasources, err := GetDatasources()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		var bodies []string
		for _, ds := range datasources {
			var err error
			var body string
			client := http.Client{Timeout: 2 * time.Second}
			resp, err := client.Get(fmt.Sprintf("http://%s:%d/api/v1/targets?state=active", ds.Host, ds.Port))
			if err != nil {
				body = `{"status":"error","errorType":"timeout"}`
			} else {
				bodyBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					body = `{"status":"error","errorType":"internal"}`
				} else {
					body = string(bodyBytes)
				}
			}
			bodies = append(bodies, body)
		}
		c.JSON(200, bodies)
	})
}

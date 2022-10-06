package server

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func routesAPIPrometheus(api *gin.RouterGroup) {

	api.GET("/prometheus/time", func(c *gin.Context) {
		log.Println("/prometheus/time")
		resp, err := http.Get("http://prometheus.kube-system:9090/api/v1/query?query=time()")
		if err != nil {
			c.JSON(500, err)
		}
		body, _ := io.ReadAll(resp.Body)
		c.String(200, string(body))
	})

	api.GET("/prometheus/namespaces", func(c *gin.Context) {
		log.Println("/prometheus/namespaces")
		resp, err := http.Get("http://prometheus.kube-system:9090/api/v1/query?query=kube_namespace_created")
		if err != nil {
			c.JSON(500, err)
		}
		body, _ := io.ReadAll(resp.Body)
		c.String(200, string(body))
	})

	api.GET("/prometheus/pods/:namespace", func(c *gin.Context) {
		namespace := c.Param("namespace")
		log.Println("/prometheus/namespaces")
		resp, err := http.Get(fmt.Sprintf(`http://prometheus.kube-system:9090/api/v1/query?query=kube_pod_created{namespace="%s"}`, namespace))
		if err != nil {
			c.JSON(500, err)
		}
		body, _ := io.ReadAll(resp.Body)
		c.String(200, string(body))
	})

	api.GET("/prometheus/nodes", func(c *gin.Context) {
		log.Println("/prometheus/nodes")
		resp, err := http.Get("http://prometheus.kube-system:9090/api/v1/query?query=kube_node_created")
		if err != nil {
			c.JSON(500, err)
		}
		body, _ := io.ReadAll(resp.Body)
		c.String(200, string(body))
	})

	api.GET("/prometheus/metadata", func(c *gin.Context) {
		log.Println("/prometheus/metedata")
		resp, err := http.Get("http://prometheus.kube-system:9090/api/v1/metadata")
		if err != nil {
			c.JSON(500, err)
		}
		body, _ := io.ReadAll(resp.Body)
		c.String(200, string(body))
	})

	api.GET("/prometheus/query", func(c *gin.Context) {
		var httpQuery HTTPQuery
		if c.ShouldBind(&httpQuery) != nil {
			c.String(400, `{"message":"invalid prometheus query request"}`)
			return
		}
		log.Println("/prometheus/query", httpQuery)
		response, err := RunHTTPPrometheusQuery(httpQuery)
		if err != nil {
			c.JSON(500, err)
		}
		c.String(200, response)
	})

	api.GET("/prometheus/query_range", func(c *gin.Context) {
		var httpQueryRange HTTPQueryRange
		if c.ShouldBind(&httpQueryRange) != nil {
			c.String(400, `{"message":"invalid prometheus query_range request"}`)
			return
		}
		log.Println("/prometheus/query_range", httpQueryRange)
		response, err := RunHTTPPrometheusQueryRange(httpQueryRange)
		if err != nil {
			c.JSON(500, err)
		}
		c.String(200, response)
	})
}

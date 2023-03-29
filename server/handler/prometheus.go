<<<<<<<< HEAD:server/api/routes_api_prometheus.go
package api
========
package handler
>>>>>>>> e244785 (baseline for reconstruction):server/handler/prometheus.go

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/server"
)

func routesAPIPrometheus(api *gin.RouterGroup) {

	api.GET("/prometheus/time", func(c *gin.Context) {
		result, err := server.InstantQuery{
			DatasourceType: server.DatasourceTypePrometheus,
			Expr:           "time()",
		}.execute()
		if err != nil {
			c.JSON(500, err)
		}
		c.String(200, result)
	})

	api.GET("/prometheus/namespaces", func(c *gin.Context) {
		result, err := server.InstantQuery{
			DatasourceType: server.DatasourceTypePrometheus,
			Expr:           "kube_namespace_created",
		}.execute()
		if err != nil {
			c.JSON(500, err)
		}
		c.String(200, result)
	})

	api.GET("/prometheus/pods/:namespace", func(c *gin.Context) {
		result, err := server.InstantQuery{
			DatasourceType: server.DatasourceTypePrometheus,
			Expr:           fmt.Sprintf(`kube_pod_created{namespace="%s"}`, c.Param("namespace")),
		}.execute()
		if err != nil {
			c.JSON(500, err)
		}
		c.String(200, result)
	})

	api.GET("/prometheus/nodes", func(c *gin.Context) {
		result, err := server.InstantQuery{
			DatasourceType: server.DatasourceTypePrometheus,
			Expr:           "kube_node_created",
		}.execute()
		if err != nil {
			c.JSON(500, err)
		}
		c.String(200, result)
	})

	api.GET("/prometheus/metadata", func(c *gin.Context) {
		result, err := server.PathQuery{
			DatasourceType: server.DatasourceTypePrometheus,
<<<<<<<< HEAD:server/api/routes_api_prometheus.go
			Path:           "/api/v1/metadata",
========
			Path:           "/handler/v1/metadata",
>>>>>>>> e244785 (baseline for reconstruction):server/handler/prometheus.go
		}.execute()
		if err != nil {
			c.JSON(500, err)
		}
		c.String(200, result)
	})

	api.GET("/prometheus/query", func(c *gin.Context) {
		result, err := server.InstantQuery{
			DatasourceType: server.DatasourceTypePrometheus,
			Expr:           c.Query("query"),
			Time:           c.Query("time"),
		}.execute()
		if err != nil {
			c.JSON(500, err)
		}
		c.String(200, result)
	})

	api.GET("/prometheus/query_range", func(c *gin.Context) {
		result, err := server.RangeQuery{
			DatasourceType: server.DatasourceTypePrometheus,
			Expr:           c.Query("query"),
			Start:          c.Query("start"),
			End:            c.Query("end"),
			Step:           c.Query("step"),
		}.execute()
		if err != nil {
			c.JSON(500, err)
		}
		c.String(200, result)
	})
}

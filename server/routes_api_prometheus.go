package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func routesAPIPrometheus(api *gin.RouterGroup) {

	api.GET("/prometheus/time", func(c *gin.Context) {
		result, err := InstantQuery{
			DatasourceType: DatasourceTypePrometheus,
			Expr:           "time()",
		}.execute()
		if err != nil {
			c.JSON(500, err)
		}
		c.String(200, result)
	})

	api.GET("/prometheus/namespaces", func(c *gin.Context) {
		result, err := InstantQuery{
			DatasourceType: DatasourceTypePrometheus,
			Expr:           "kube_namespace_created",
		}.execute()
		if err != nil {
			c.JSON(500, err)
		}
		c.String(200, result)
	})

	api.GET("/prometheus/pods/:namespace", func(c *gin.Context) {
		result, err := InstantQuery{
			DatasourceType: DatasourceTypePrometheus,
			Expr:           fmt.Sprintf(`kube_pod_created{namespace="%s"}`, c.Param("namespace")),
		}.execute()
		if err != nil {
			c.JSON(500, err)
		}
		c.String(200, result)
	})

	api.GET("/prometheus/nodes", func(c *gin.Context) {
		result, err := InstantQuery{
			DatasourceType: DatasourceTypePrometheus,
			Expr:           "kube_node_created",
		}.execute()
		if err != nil {
			c.JSON(500, err)
		}
		c.String(200, result)
	})

	api.GET("/prometheus/metadata", func(c *gin.Context) {
		result, err := PathQuery{
			DatasourceType: DatasourceTypePrometheus,
			Path:           "/api/v1/metadata",
		}.execute()
		if err != nil {
			c.JSON(500, err)
		}
		c.String(200, result)
	})

	api.GET("/prometheus/query", func(c *gin.Context) {
		result, err := InstantQuery{
			DatasourceType: DatasourceTypePrometheus,
			Expr:           c.Query("query"),
			Time:           c.Query("time"),
		}.execute()
		if err != nil {
			c.JSON(500, err)
		}
		c.String(200, result)
	})

	api.GET("/prometheus/query_range", func(c *gin.Context) {
		result, err := RangeQuery{
			DatasourceType: DatasourceTypePrometheus,
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

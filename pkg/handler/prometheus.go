package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type prometheusHandler struct {
	// todo : something able to query
}

// GET /prometheus/metadata
func (ph *prometheusHandler) nodes(c *gin.Context) {
	result, err := server.PathQuery{
		DatasourceType: server.DatasourceTypePrometheus,
		Path:           "/api/v1/metadata",
	}.execute()
	if err != nil {
		c.JSON(500, err)
	}
	c.String(200, result)
}

// GET /prometheus/time
func (ph *prometheusHandler) time(c *gin.Context) {
	result, err := server.InstantQuery{
		DatasourceType: server.DatasourceTypePrometheus,
		Expr:           "time()",
	}.execute()
	if err != nil {
		c.JSON(500, err)
	}
	c.String(200, result)
}

// GET /prometheus/namespaces
func (ph *prometheusHandler) namespace(c *gin.Context) {
	result, err := server.InstantQuery{
		DatasourceType: server.DatasourceTypePrometheus,
		Expr:           "kube_namespace_created",
	}.execute()
	if err != nil {
		c.JSON(500, err)
	}
	c.String(200, result)
}

// GET /prometheus/pods/:namespace
func (ph *prometheusHandler) pod(c *gin.Context) {
	result, err := server.InstantQuery{
		DatasourceType: server.DatasourceTypePrometheus,
		Expr:           fmt.Sprintf(`kube_pod_created{namespace="%s"}`, c.Param("namespace")),
	}.execute()
	if err != nil {
		c.JSON(500, err)
	}
	c.String(200, result)
}

// GET /prometheus/nodes
func (ph *prometheusHandler) nodes(c *gin.Context) {
	result, err := server.InstantQuery{
		DatasourceType: server.DatasourceTypePrometheus,
		Expr:           "kube_node_created",
	}.execute()
	if err != nil {
		c.JSON(500, err)
	}
	c.String(200, result)
}

// GET /prometheus/query
func (ph *prometheusHandler) query(c *gin.Context) {
	result, err := server.InstantQuery{
		DatasourceType: server.DatasourceTypePrometheus,
		Expr:           c.Query("query"),
		Time:           c.Query("time"),
	}.execute()
	if err != nil {
		c.JSON(500, err)
	}
	c.String(200, result)
}

// GET /prometheus/query_range
func (ph *prometheusHandler) queryRange(c *gin.Context) {

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
}

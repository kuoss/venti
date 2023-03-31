package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/query"
)

type prometheusHandler struct {
	querier query.Querier
}

func NewPrometheusHandler(querier query.Querier) *prometheusHandler {
	return &prometheusHandler{
		querier: querier,
	}
}

// GET /prometheus/metadata
func (ph *prometheusHandler) Metadata(c *gin.Context) {

	qr, err := ph.querier.Execute(c.Request.Context(), query.Query{
		Path: "metadata",
	})
	if err != nil {
		c.JSON(500, err)
	}

	c.String(200, qr)
}

// GET /prometheus/query
func (ph *prometheusHandler) Query(c *gin.Context) {

	qr, err := ph.querier.Execute(c.Request.Context(), query.Query{
		Path: "/api/v1/query",
		Params: map[string]string{
			"query": c.Query("query"),
			"time":  c.Query("time"),
		},
	})
	if err != nil {
		c.JSON(500, err)
	}
	c.String(200, qr)
}

// GET /prometheus/query_range
func (ph *prometheusHandler) QueryRange(c *gin.Context) {

	qr, err := ph.querier.Execute(c.Request.Context(), query.Query{
		Path: "/api/v1/query",
		Params: map[string]string{
			"query": c.Query("query"),
			"time":  c.Query("time"),
			"start": c.Query("start"),
			"end":   c.Query("end"),
			"step":  c.Query("step"),
		},
	})
	if err != nil {
		c.JSON(500, err)
	}
	c.String(200, qr)
}

// this is also another isntantQuery
// GET /prometheus/time -> /prometheus/query?expr=time()
// todo: migrate with /prometheus/query
func (ph *prometheusHandler) time(c *gin.Context) {

	/*
		result, err := server.InstantQuery{
			DatasourceType: server.DatasourceTypePrometheus,
			Expr:           "time()",
		}.execute()
		if err != nil {
			c.JSON(500, err)
		}
			c.String(200, result)
	*/
}

// GET /prometheus/namespaces
// todo: migrate with /prometheus/query
func (ph *prometheusHandler) namespaces(c *gin.Context) {
	/*
		result, err := server.InstantQuery{
			DatasourceType: server.DatasourceTypePrometheus,
			Expr:           "kube_namespace_created",
		}.execute()
		if err != nil {
			c.JSON(500, err)
		}
		c.String(200, result)
	*/
}

// GET /prometheus/pods/:namespace
// todo: migrate with /prometheus/query?kube_pod_created{namespace="namespace01"}
func (ph *prometheusHandler) pod(c *gin.Context) {
	/*
		result, err := server.InstantQuery{
			DatasourceType: server.DatasourceTypePrometheus,
			Expr:           fmt.Sprintf(`kube_pod_created{namespace="%s"}`, c.Param("namespace")),
		}.execute()
		if err != nil {
			c.JSON(500, err)
		}
		c.String(200, result)
	*/
}

// GET /prometheus/nodes
// todo: migrate with /prometheus/query?kube_node_created
func (ph *prometheusHandler) nodes(c *gin.Context) {
	/*
		result, err := server.InstantQuery{
			DatasourceType: server.DatasourceTypePrometheus,
			Expr:           "kube_node_created",
		}.execute()
		if err != nil {
			c.JSON(500, err)
		}
		c.String(200, result)
	*/
}

package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg/query"
)

type letheHandler struct {
	// todo: something to query lethe data
	querier query.Querier
}

func NewLetheHandler(querier query.Querier) *letheHandler {
	return &letheHandler{
		querier: querier,
	}
}

// GET /lethe/metadata
func (lh *letheHandler) Metadata(c *gin.Context) {
	qr, err := lh.querier.Execute(c.Request.Context(), query.Query{
		Path: "metadata",
	})
	if err != nil {
		c.JSON(500, err)
	}

	c.String(200, qr)
}

// GET /lethe/query
func (lh *letheHandler) Query(c *gin.Context) {
	qr, err := lh.querier.Execute(c.Request.Context(), query.Query{
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

// GET /lethe/query_range
func (lh *letheHandler) QueryRange(c *gin.Context) {
	qr, err := lh.querier.Execute(c.Request.Context(), query.Query{
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

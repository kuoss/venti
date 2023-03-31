package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/pkg"
)

type letheHandler struct {
	// todo: something to query lethe data
}

// GET /lethe/metadata
func (lh *letheHandler) metadata(c *gin.Context) {
	result, err := pkg.PathQuery{
		DatasourceType: pkg.DatasourceTypeLethe,
		Path:           "/api/v1/metadata",
	}.execute()
	if err != nil {
		c.JSON(500, err)
	}
	c.String(200, result)
}

// GET /lethe/query
func (lh *letheHandler) query(c *gin.Context) {
	result, err := pkg.InstantQuery{
		DatasourceType: pkg.DatasourceTypeLethe,
		Expr:           c.Query("query"),
		Time:           c.Query("time"),
	}.execute()
	if err != nil {
		c.JSON(500, err)
	}
	c.String(200, result)
}

// GET /lethe/query_range
func (lh *letheHandler) query_range(c *gin.Context) {
	result, err := pkg.RangeQuery{
		DatasourceType: pkg.DatasourceTypeLethe,
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

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kuoss/venti/server"
)

func routesAPILethe(api *gin.RouterGroup) {
	api.GET("/lethe/metadata", func(c *gin.Context) {
		result, err := server.PathQuery{
			DatasourceType: server.DatasourceTypeLethe,
			Path:           "/api/v1/metadata",
		}.execute()
		if err != nil {
			c.JSON(500, err)
		}
		c.String(200, result)
	})

	api.GET("/lethe/query", func(c *gin.Context) {
		result, err := server.InstantQuery{
			DatasourceType: server.DatasourceTypeLethe,
			Expr:           c.Query("query"),
			Time:           c.Query("time"),
		}.execute()
		if err != nil {
			c.JSON(500, err)
		}
		c.String(200, result)
	})

	api.GET("/lethe/query_range", func(c *gin.Context) {
		result, err := server.RangeQuery{
			DatasourceType: server.DatasourceTypeLethe,
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

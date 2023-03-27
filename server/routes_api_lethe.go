package server

import (
	"github.com/gin-gonic/gin"
)

func routesAPILethe(api *gin.RouterGroup) {

	api.GET("/lethe/metadata", func(c *gin.Context) {
		result, err := PathQuery{
			DatasourceType: DatasourceTypeLethe,
			Path:           "/api/v1/metadata",
		}.execute()
		if err != nil {
			c.JSON(500, err)
		}
		c.String(200, result)
	})

	api.GET("/lethe/query", func(c *gin.Context) {
		result, err := InstantQuery{
			DatasourceType: DatasourceTypeLethe,
			Expr:           c.Query("query"),
			Time:           c.Query("time"),
		}.execute()
		if err != nil {
			c.JSON(500, err)
		}
		c.String(200, result)
	})

	api.GET("/lethe/query_range", func(c *gin.Context) {
		result, err := RangeQuery{
			DatasourceType: DatasourceTypeLethe,
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

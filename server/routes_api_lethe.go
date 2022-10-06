package server

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func routesAPILethe(api *gin.RouterGroup) {

	api.GET("/lethe/metadata", func(c *gin.Context) {
		log.Println("/lethe/metadata")
		req, _ := http.NewRequest("GET", "http://lethe:8080/api/v1/metadata", nil)
		resp, _ := http.Get(req.URL.String())
		body, _ := io.ReadAll(resp.Body)
		c.String(200, string(body))
	})

	api.GET("/lethe/query", func(c *gin.Context) {
		var httpQuery HTTPQuery
		if c.ShouldBind(&httpQuery) != nil {
			c.String(400, `{"message":"invalid lethe query request"}`)
			return
		}
		log.Println("/lethe/query", httpQuery)
		response, err := RunHTTPLetheQuery(httpQuery)
		if err != nil {
			c.JSON(500, err)
		}
		c.String(200, response)
	})

	api.GET("/lethe/query_range", func(c *gin.Context) {
		var httpQueryRange HTTPQueryRange
		if c.ShouldBind(&httpQueryRange) != nil {
			c.String(400, `{"message":"invalid lethe query_range request"}`)
			return
		}
		log.Println("/lethe/query_range", httpQueryRange)
		response, err := RunHTTPLetheQueryRange(httpQueryRange)
		if err != nil {
			c.JSON(500, err)
		}
		c.String(200, response)
	})
}

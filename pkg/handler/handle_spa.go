package handler

import (
	"net/http"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func handleSPA() gin.HandlerFunc {
	directory := static.LocalFile("./web/dist", true)
	fileserver := http.StripPrefix("/", http.FileServer(directory))
	return func(c *gin.Context) {
		if !directory.Exists("/", c.Request.URL.Path) {
			c.Request.URL.Path = "/"
		}
		fileserver.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}

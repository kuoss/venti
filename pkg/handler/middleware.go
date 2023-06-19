package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func tokenRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "valid token required",
			})
			return
		}
		at := strings.Split(token, " ")
		if len(at) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "valid token required",
			})
			return
		}
		token = at[1]
		// fixme: decide token
		if token != "fixme" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "valid token required",
			})
			return
		}
		c.Next()
	}
}

package handler

import (
	"fmt"
	"github.com/kuoss/venti/pkg/handler/api"
	"strings"

	"github.com/gin-gonic/gin"
)

func tokenRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			api.ResponseError(c, api.ErrorUnauthorized, fmt.Errorf("token required"))
			c.Abort()
		}
		at := strings.Split(token, " ")
		if len(at) != 2 {
			api.ResponseError(c, api.ErrorUnauthorized, fmt.Errorf("valid token required"))
			c.Abort()
		}
		token = at[1]
		// fixme: decide token
		if token != "fixme" {
			api.ResponseError(c, api.ErrorUnauthorized, fmt.Errorf("valid token required"))
			c.Abort()
		}
		c.Next()
	}
}

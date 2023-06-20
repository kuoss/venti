package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestInvalidToken(t *testing.T) {

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/test", tokenRequired(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "test",
		})
	})

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer INVALID")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, "{\"error\":\"valid token required\",\"errorType\":\"unauthorized\",\"status\":\"error\"}", w.Body.String())
}

func TestValidToken(t *testing.T) {

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/test", tokenRequired(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "test",
		})
	})

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/test", nil)

	req.Header.Set("Authorization", "Bearer fixme")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"message\":\"test\"}", w.Body.String())
}

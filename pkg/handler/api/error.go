package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResponseError(c *gin.Context, typ errorType, err error) {
	// https://prometheus.io/docs/prometheus/latest/querying/api/#format-overview
	c.JSON(getCodeFromType(typ), gin.H{
		"status":    "error",
		"errorType": typ,
		"error":     err.Error(),
	})
}

// https://github.com/prometheus/prometheus/blob/main/web/api/v1/api.go#L72
type errorType string

const (
	ErrorCanceled     errorType = "canceled"     //
	ErrorExec         errorType = "execution"    //
	ErrorUnauthorized errorType = "unauthorized" // 401 Unauthorized
	ErrorNotFound     errorType = "not_found"    // 404 Not Found
	ErrorBadData      errorType = "bad_data"     // 405 StatusMethodNotAllowed
	ErrorTimeout      errorType = "timeout"      // 408 Request Timeout
	ErrorInternal     errorType = "internal"     // 500 Internal Server Error
	ErrorUnavailable  errorType = "unavailable"  // 503 Service Unavailable
)

type Error struct {
	Type errorType
	Err  error
}

func getCodeFromType(typ errorType) int {
	switch typ {
	case ErrorUnauthorized:
		return http.StatusUnauthorized // 401 Unauthorized
	case ErrorNotFound:
		return http.StatusNotFound // 404 Not Found
	case ErrorBadData:
		return http.StatusMethodNotAllowed // 405 StatusMethodNotAllowed
	case ErrorTimeout:
		return http.StatusRequestTimeout // 408 Request Timeout
	case ErrorInternal:
		return http.StatusInternalServerError // 500 Internal Server Error
	case ErrorUnavailable:
		return http.StatusServiceUnavailable // 503 Service Unavailable
	}
	return http.StatusServiceUnavailable // 503 Service Unavailable
}

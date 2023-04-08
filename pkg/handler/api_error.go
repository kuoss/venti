package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func responseAPIError(c *gin.Context, apiError *apiError) {
	// https://prometheus.io/docs/prometheus/latest/querying/api/#format-overview
	var code int
	switch apiError.typ {
	case errorBadData:
		code = http.StatusBadRequest
	default:
		code = http.StatusServiceUnavailable
	}
	c.JSON(code, gin.H{
		"status":    "error",
		"errorType": apiError.typ,
		"error":     apiError.Error(),
	})
}

// https://github.com/prometheus/prometheus/blob/main/web/api/v1/api.go#L72
type errorType string

const (
	errorNone        errorType = ""
	errorTimeout     errorType = "timeout"
	errorCanceled    errorType = "canceled"
	errorExec        errorType = "execution"
	errorBadData     errorType = "bad_data"
	errorInternal    errorType = "internal"
	errorUnavailable errorType = "unavailable"
	errorNotFound    errorType = "not_found"
)

type apiError struct {
	typ errorType
	err error
}

func (e *apiError) Error() string {
	return fmt.Sprintf("%s: %s", e.typ, e.err)
}

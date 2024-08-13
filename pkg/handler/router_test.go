package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	assert.NotEmpty(t, handlers)
	assert.NotEmpty(t, handlers.alertHandler)
	assert.NotEmpty(t, handlers.authHandler)
	assert.NotEmpty(t, handlers.dashboardHandler)
	assert.NotEmpty(t, handlers.datasourceHandler)
	assert.NotEmpty(t, handlers.remoteHandler)
	assert.NotEmpty(t, handlers.statusHandler)

	router := NewRouter(services)
	assert.NotEmpty(t, router)
}

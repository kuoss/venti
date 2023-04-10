package handler

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadHandlers(t *testing.T) {
	assert.NotNil(t, handlers)
	assert.NotNil(t, handlers.alertHandler)
	assert.NotNil(t, handlers.authHandler)
	assert.NotNil(t, handlers.configHandler)
	assert.NotNil(t, handlers.dashboardHandler)
	assert.NotNil(t, handlers.datasourceHandler)
	assert.NotNil(t, handlers.remoteHandler)
}
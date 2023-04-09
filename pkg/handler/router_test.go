package handler

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetupRouter(t *testing.T) {
	assert.NotNil(t, router)
}

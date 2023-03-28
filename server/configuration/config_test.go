package configuration

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	Load("Unknown")
}

func TestGetConfig(t *testing.T) {
	config := GetConfig()
	assert.Equal(t, "Unknown", config.Version)
}

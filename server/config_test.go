package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	LoadConfig("Unknown")
}

func TestGetConfig(t *testing.T) {
	config := GetConfig()
	assert.Equal(t, "Unknown", config.Version)
}


func TestListServices(t *testing.T) {
	services, err := listServices()
}
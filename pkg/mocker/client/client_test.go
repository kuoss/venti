package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	endpoint := "hello"
	client := New(endpoint)
	assert.Equal(t, endpoint, client.endpoint)
}

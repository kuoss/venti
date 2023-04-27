package alertmanager_test

import (
	"testing"

	"github.com/kuoss/venti/pkg/mocker"
	"github.com/kuoss/venti/pkg/mocker/alertmanager"
	mockerClient "github.com/kuoss/venti/pkg/mocker/client"
	"github.com/stretchr/testify/assert"
)

var (
	server *mocker.Server
	client *mockerClient.Client
)

func init() {
	server, _ = alertmanager.New(0)
	client = mockerClient.New(server.URL)
}

func Test_api_v1_alerts(t *testing.T) {
	code, body, err := client.GET("/api/v1/alerts", "")
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.JSONEq(t, `{"status":"success","data":{}}`, body)
}

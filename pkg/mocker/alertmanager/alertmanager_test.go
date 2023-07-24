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
	var err error
	server, err = alertmanager.New()
	if err != nil {
		panic(err)
	}
	client = mockerClient.New(server.URL)
}

func Test_api_v1_status_buildinfo(t *testing.T) {
	code, body, err := client.GET("/api/v1/status/buildinfo", "")
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.JSONEq(t, `{"status":"success","data":{"version":"2.41.0-alertmanager","revision":"c0d8a56c69014279464c0e15d8bfb0e153af0dab","branch":"HEAD","buildUser":"root@d20a03e77067","buildDate":"20221220-10:40:45","goVersion":"go1.19.4"}}`, body)
}

func Test_api_v1_alerts(t *testing.T) {
	code, body, err := client.GET("/api/v1/alerts", "")
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.JSONEq(t, `{"status":"success","data":{}}`, body)
}

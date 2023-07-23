package servers_test

import (
	"testing"

	ms "github.com/kuoss/venti/pkg/mock/servers"
	"github.com/stretchr/testify/assert"
)

var servers *ms.Servers

func init() {
	servers = ms.New(ms.Requirements{
		{Type: ms.TypeAlertmanager, Name: "alertmanager1", IsMain: false, BasicAuth: false},
		{Type: ms.TypeLethe, Name: "lethe1", IsMain: true, BasicAuth: false},
		{Type: ms.TypeLethe, Name: "lethe2", IsMain: false, BasicAuth: false},
		{Type: ms.TypePrometheus, Name: "prometheus1", IsMain: true, BasicAuth: false},
		{Type: ms.TypePrometheus, Name: "prometheus2", IsMain: false, BasicAuth: false},
		{Type: ms.TypePrometheus, Name: "prometheus3", IsMain: false, BasicAuth: true},
	})
}

func TestNew(t *testing.T) {
	assert.Equal(t, 6, len(servers.Svrs))
	assert.Equal(t, "alertmanager1", servers.Svrs[0].Name)
	assert.Equal(t, "prometheus3", servers.Svrs[5].Name)
}

func TestGetDatasources(t *testing.T) {
	assert.Equal(t, 5, len(servers.GetDatasources())) // alertmanager cannot be a datasource
}

func GetServersByType(t *testing.T) {
	assert.Equal(t, 1, len(servers.GetServersByType(ms.TypeAlertmanager)))
	assert.Equal(t, 2, len(servers.GetServersByType(ms.TypeLethe)))
	assert.Equal(t, 3, len(servers.GetServersByType(ms.TypePrometheus)))
}

func TestClose(t *testing.T) {
	servers.Close()
}

package mocker_test

import (
	"net/url"
	"testing"

	"github.com/kuoss/venti/pkg/mocker"
	mockerClient "github.com/kuoss/venti/pkg/mocker/client"
	"github.com/stretchr/testify/assert"
)

var (
	server *mocker.Server
	client *mockerClient.Client
)

func init() {
	server = mocker.New()
	server.GET("/ping", func(c *mocker.Context) {
		c.JSON(200, mocker.H{"message": "pong"})
	})
	server.GET("/ping2", func(c *mocker.Context) {
		c.JSONString(200, `{"message":"pong"}`)
	})
	server.GET("/greet", func(c *mocker.Context) {
		name := c.Query("name")
		c.JSON(200, mocker.H{
			"message": "Hello_" + name,
		})
	})
	_ = server.Start(0)

	client = mockerClient.New(server.URL)
}

func TestNew_0(t *testing.T) {
	assert.NotZero(t, server)
	u, _ := url.Parse(server.URL)
	assert.Equal(t, "127.0.0.1", u.Hostname())
	assert.NotEqual(t, "0", u.Port())
}

func TestNew_9999(t *testing.T) {
	tempServer := mocker.New()
	err := tempServer.Start(9999)
	assert.Nil(t, err)
	defer tempServer.Close()
	assert.Equal(t, "http://127.0.0.1:9999", tempServer.URL)
}

func TestNew_dup(t *testing.T) {
	var err error

	tempServer1 := mocker.New()
	err = tempServer1.Start(9999)
	assert.Nil(t, err)
	defer tempServer1.Close()

	tempServer2 := mocker.New()
	err = tempServer2.Start(9999)
	assert.Equal(t, "error on Listen: listen tcp 127.0.0.1:9999: bind: address already in use", err.Error())
}

func TestNew_99999(t *testing.T) {
	tempServer := mocker.New()
	err := tempServer.Start(99999)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "invalid port number: 99999")
}

func TestNew_basicAuth(t *testing.T) {
	tempServer := mocker.New()
	tempServer.SetBasicAuth("abc", "123")
	tempServer.GET("/ping", func(c *mocker.Context) {
		c.JSON(200, mocker.H{"message": "pong"})
	})
	_ = tempServer.Start(0)
	tempClient := mockerClient.New(tempServer.URL)
	code, body, err := tempClient.GET("/ping", "")
	assert.NoError(t, err)
	assert.Equal(t, 401, code)
	assert.Equal(t, "401 unauthorized\n", body)
}

func Test_ping(t *testing.T) {
	code, body, err := client.GET("/ping", "")
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.JSONEq(t, `{"message":"pong"}`, body)
}

func Test_ping2(t *testing.T) {
	code, body, err := client.GET("/ping2", "")
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.JSONEq(t, `{"message":"pong"}`, body)
}

func Test_ping3(t *testing.T) {
	code, body, err := client.GET("/ping2", "")
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.JSONEq(t, `{"message":"pong"}`, body)
}

func Test_greet_empty(t *testing.T) {
	code, body, err := client.GET("/greet", "")
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.JSONEq(t, `{"message":"Hello_"}`, body)
}

func Test_greet_John(t *testing.T) {
	code, body, err := client.GET("/greet", "name=John")
	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.JSONEq(t, `{"message":"Hello_John"}`, body)
}

func TestGET_404_page_not_found(t *testing.T) {
	code, body, err := client.GET("/not_found", "")
	assert.NoError(t, err)
	assert.Equal(t, 404, code)
	assert.Equal(t, "404 page not found\n", body)
}

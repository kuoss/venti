package mocker_test

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/kuoss/venti/pkg/mocker"
	mockerClient "github.com/kuoss/venti/pkg/mocker/client"
	"github.com/stretchr/testify/assert"
)

var (
	mainServer *mocker.Server
	mainClient *mockerClient.Client
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	mainServer = mocker.New()
	mainServer.GET("/ping", func(c *mocker.Context) {
		c.JSON(200, mocker.H{"message": "pong"})
	})
	mainServer.GET("/ping2", func(c *mocker.Context) {
		c.JSONString(200, `{"message":"pong"}`)
	})
	mainServer.GET("/greet", func(c *mocker.Context) {
		name := c.Query("name")
		c.JSON(200, mocker.H{
			"message": "Hello_" + name,
		})
	})
	err := mainServer.Start(0)
	if err != nil {
		panic(err)
	}
	mainClient = mockerClient.New(mainServer.URL)
}

func shutdown() {
	mainServer.Close()
}

func TestNew(t *testing.T) {
	// default case
	assert.NotZero(t, mainServer)
	u, err := url.Parse(mainServer.URL)
	assert.NoError(t, err)
	assert.Equal(t, "127.0.0.1", u.Hostname())
	assert.NotEqual(t, "0", u.Port())
}

func TestStart(t *testing.T) {
	testCases := []struct {
		port      int
		wantError string
		wantURL   string
	}{
		{
			9999,
			"",
			"http://127.0.0.1:9999",
		},
		{
			99999,
			"invalid port number: 99999",
			"",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(tt *testing.T) {
			tempServer := mocker.New()
			err := tempServer.Start(tc.port)
			if tc.wantError == "" {
				defer tempServer.Close()
				assert.NoError(tt, err)
			} else {
				assert.EqualError(tt, err, tc.wantError)
			}
			assert.Equal(tt, tc.wantURL, tempServer.URL)
		})
	}
}

func TestNew_duplicate(t *testing.T) {
	var err error

	tempServer1 := mocker.New()
	err = tempServer1.Start(1234)
	assert.Nil(t, err)
	defer tempServer1.Close()

	tempServer2 := mocker.New()
	err = tempServer2.Start(1234)
	assert.Equal(t, "error on Listen: listen tcp 127.0.0.1:1234: bind: address already in use", err.Error())
}

func TestSetBasicAuth(t *testing.T) {
	tempServer := mocker.New()
	tempServer.SetBasicAuth("abc", "123")
	tempServer.GET("/ping", func(c *mocker.Context) {
		c.JSON(200, mocker.H{"message": "pong"})
	})
	err := tempServer.Start(0)
	assert.NoError(t, err)

	tempClient1 := mockerClient.New(tempServer.URL)

	tempClient2 := mockerClient.New(tempServer.URL)
	tempClient2.SetBasicAuth("xxx", "xxx")

	tempClient3 := mockerClient.New(tempServer.URL)
	tempClient3.SetBasicAuth("abc", "123")

	testCases := []struct {
		client   *mockerClient.Client
		wantCode int
		wantBody string
	}{
		{
			tempClient1,
			401, "401 unauthorized\n",
		},
		{
			tempClient2,
			401, "401 unauthorized\n",
		},
		{
			tempClient3,
			200, `{"message":"pong"}`,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(tt *testing.T) {
			code, body, err := tc.client.GET("/ping", "")
			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, code)
			assert.Equal(t, tc.wantBody, body)
		})
	}
}

func Test_ping(t *testing.T) {
	testCases := []struct {
		path      string
		rawQuery  string
		wantCode  int
		wantBody  string
		wantError string
	}{
		{
			"/ping", "",
			200, `{"message":"pong"}`, "",
		},
		{
			"/ping2", "",
			200, `{"message":"pong"}`, "",
		},
		{
			"/greet", "",
			200, `{"message":"Hello_"}`, "",
		},
		{
			"/greet", "name=John",
			200, `{"message":"Hello_John"}`, "",
		},
		{
			"/not_found", "",
			404, "404 page not found\n", "",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(tt *testing.T) {
			code, body, err := mainClient.GET(tc.path, tc.rawQuery)
			if tc.wantError == "" {
				assert.NoError(tt, err)
			} else {
				assert.EqualError(tt, err, tc.wantError)
			}
			assert.Equal(tt, tc.wantCode, code)
			assert.Equal(tt, tc.wantBody, body)
		})
	}
}

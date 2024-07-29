package mocker_test

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/kuoss/venti/pkg/mocker"
	mockerClient "github.com/kuoss/venti/pkg/mocker/client"
	"github.com/stretchr/testify/require"
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
	err := mainServer.Start()
	if err != nil {
		panic(err)
	}
	mainClient = mockerClient.New(mainServer.URL())
}

func shutdown() {
	mainServer.Close()
}

func TestNew(t *testing.T) {
	require.NotZero(t, mainServer)
	u, err := url.Parse(mainServer.URL())
	require.NoError(t, err)
	require.Equal(t, "http", u.Scheme)
	require.Equal(t, "127.0.0.1", u.Hostname())
}

func TestStart(t *testing.T) {
	tempServer := mocker.NewWithPort(12345)
	err := tempServer.Start()
	require.NoError(t, err)
}

func TestClose(t *testing.T) {
	tempServer := mocker.NewWithPort(12346)
	err := tempServer.Start()
	require.NoError(t, err)
	tempServer.Close()
}

func TestSetBasicAuth(t *testing.T) {
	tempServer := mocker.New()
	tempServer.SetBasicAuth("abc", "123")
	tempServer.GET("/ping", func(c *mocker.Context) {
		c.JSON(200, mocker.H{"message": "pong"})
	})
	err := tempServer.Start()
	require.NoError(t, err)

	tempClient1 := mockerClient.New(tempServer.URL())

	tempClient2 := mockerClient.New(tempServer.URL())
	tempClient2.SetBasicAuth("xxx", "xxx")

	tempClient3 := mockerClient.New(tempServer.URL())
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
			require.NoError(t, err)
			require.Equal(t, tc.wantCode, code)
			require.Equal(t, tc.wantBody, body)
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
				require.NoError(tt, err)
			} else {
				require.EqualError(tt, err, tc.wantError)
			}
			require.Equal(tt, tc.wantCode, code)
			require.Equal(tt, tc.wantBody, body)
		})
	}
}

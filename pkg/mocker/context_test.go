package mocker

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin/render"
	"github.com/stretchr/testify/require"
)

func TestJSON(t *testing.T) {
	testCases := []struct {
		code     int
		obj      any
		wantCode int
		wantBody string
	}{
		{
			200, "hello",
			200, `"hello"`,
		},
		{
			405, "hello",
			405, `"hello"`,
		},
		{
			200, []int{1, 2, 3},
			200, `[1,2,3]`,
		},
		{
			200, map[string]string{"hello": "world", "lorem": "ipsum"},
			200, `{"hello":"world","lorem":"ipsum"}`,
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			w := httptest.NewRecorder()
			context := &Context{Writer: w, Request: &http.Request{}}
			context.JSON(tc.code, tc.obj)

			require.Equal(t, tc.wantCode, w.Code)
			require.Equal(t, tc.wantBody, w.Body.String())
		})
	}
}

func TestJSONString(t *testing.T) {
	testCases := []struct {
		code     int
		str      string
		wantCode int
		wantBody string
	}{
		{
			200, "hello",
			200, "hello",
		},
		{
			405, "hello",
			405, "hello",
		},
		{
			200, `[1,2,3]`,
			200, `[1,2,3]`,
		},
		{
			200, `{"hello":"world","lorem":"ipsum"}`,
			200, `{"hello":"world","lorem":"ipsum"}`,
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			w := httptest.NewRecorder()
			context := &Context{Writer: w, Request: &http.Request{}}
			context.JSONString(tc.code, tc.str)

			require.Equal(t, tc.wantCode, w.Code)
			require.Equal(t, tc.wantBody, w.Body.String())
		})
	}
}

func TestRender(t *testing.T) {
	testCases := []struct {
		code     int
		render   render.Render
		wantCode int
		wantBody string
	}{
		// valid data
		{
			200, render.JSON{Data: "hello"},
			200, `"hello"`,
		},
		{
			405, render.JSON{Data: "hello"},
			405, `"hello"`,
		},
		{
			200, render.JSON{Data: map[string]string{"hello": "world", "lorem": "ipsum"}},
			200, `{"hello":"world","lorem":"ipsum"}`,
		},
		{
			200, render.Data{ContentType: "application/json", Data: []byte(`{"hello":"world","lorem":"ipsum"}`)},
			200, `{"hello":"world","lorem":"ipsum"}`,
		},
		// invalid data
		{
			200, render.JSON{Data: make(chan int)},
			200, ``,
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			w := httptest.NewRecorder()
			context := &Context{Writer: w, Request: &http.Request{}}
			context.render(tc.code, tc.render)

			require.Equal(t, tc.wantCode, w.Code)
			require.Equal(t, tc.wantBody, w.Body.String())
		})
	}
}

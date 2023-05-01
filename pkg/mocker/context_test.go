package mocker

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin/render"
	"github.com/stretchr/testify/assert"
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
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			w := httptest.NewRecorder()
			context := &Context{Writer: w, Request: &http.Request{}}
			context.JSON(tc.code, tc.obj)

			assert.Equal(t, tc.wantCode, w.Code)
			assert.Equal(t, tc.wantBody, w.Body.String())
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
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			w := httptest.NewRecorder()
			context := &Context{Writer: w, Request: &http.Request{}}
			context.JSONString(tc.code, tc.str)

			assert.Equal(t, tc.wantCode, w.Code)
			assert.Equal(t, tc.wantBody, w.Body.String())
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
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			w := httptest.NewRecorder()
			context := &Context{Writer: w, Request: &http.Request{}}
			context.render(tc.code, tc.render)

			assert.Equal(t, tc.wantCode, w.Code)
			assert.Equal(t, tc.wantBody, w.Body.String())
		})
	}
}

// Currently, the library gin panics, so instead of testing for error, we test for panic. (package coverage: 97.1%)
// TODO: bump up & error test https://github.com/gin-gonic/gin/commit/0c96a20209ca035964be126a745c167196fb6db3
func TestRenderError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	w := httptest.NewRecorder()
	data := make(chan int)
	assert.Error(t, (render.JSON{Data: data}).Render(w))
}

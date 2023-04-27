package mocker

import (
	"net/http"

	"github.com/gin-gonic/gin/render"
)

// https://github.com/gin-gonic/gin/blob/v1.9.0/context.go

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
}

func (c *Context) JSON(code int, obj any) {
	c.Render(code, render.JSON{Data: obj})
}

func (c *Context) JSONString(code int, str string) {
	c.Render(code, render.Data{ContentType: "application/json", Data: []byte(str)})
}

func (c *Context) Render(code int, r render.Render) {
	c.Writer.WriteHeader(code)
	_ = r.Render(c.Writer)
}

func (c *Context) Query(key string) (value string) {
	return c.Request.URL.Query().Get(key)
}

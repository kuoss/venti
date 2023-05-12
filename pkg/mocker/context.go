package mocker

import (
	"net/http"

	"github.com/gin-gonic/gin/render"
	"github.com/kuoss/common/logger"
)

// https://github.com/gin-gonic/gin/blob/v1.9.0/context.go

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
}

func (c *Context) JSON(code int, obj any) {
	c.render(code, render.JSON{Data: obj})
}

func (c *Context) JSONString(code int, str string) {
	c.render(code, render.Data{ContentType: "application/json", Data: []byte(str)})
}

func (c *Context) render(code int, r render.Render) {
	c.Writer.WriteHeader(code)
	err := r.Render(c.Writer)
	if err != nil {
		logger.Errorf("Render err: %s", err)
	}
}

func (c *Context) Query(key string) (value string) {
	return c.Request.URL.Query().Get(key)
}

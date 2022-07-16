package gtp

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	/* origin objects */
	Writer http.ResponseWriter
	Req    *http.Request

	/* request info */
	Path   string
	Method string
	Params map[string]string // store the router parameter

	/* response info */
	StatusCode int

	/* middleware */
	handlers []HandlerFunc
	index    int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

func (ctx *Context) PostForm(key string) string {
	return ctx.Req.FormValue(key)
}

func (ctx *Context) Query(key string) string {
	return ctx.Req.URL.Query().Get(key)
}

func (ctx *Context) Param(key string) string {
	value, _ := ctx.Params[key]
	return value
}

func (ctx *Context) Status(code int) {
	ctx.StatusCode = code
	ctx.Writer.WriteHeader(code)
}

func (ctx *Context) Next() {
	ctx.index++
	s := len(ctx.handlers)
	for ; ctx.index < s; ctx.index++ {
		ctx.handlers[ctx.index](ctx)
	}
}

func (ctx *Context) Fail(code int, err string) {
	ctx.index = len(ctx.handlers)
	ctx.JSON(code, H{"message": err})
}

func (ctx *Context) Header(key, value string) {
	if value == "" {
		ctx.Writer.Header().Del(key)
		return
	}
	ctx.Writer.Header().Set(key, value)
}

func (ctx *Context) Data(code int, data []byte) {
	ctx.Status(code)
	ctx.Writer.Write(data)
}

func (ctx *Context) String(code int, format string, values ...interface{}) {
	ctx.Header("Content-Type", "text/plain")
	ctx.Status(code)
	ctx.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (ctx *Context) JSON(code int, obj interface{}) {
	ctx.Header("Content-Type", "application/json")
	ctx.Status(code)
	encoder := json.NewEncoder(ctx.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(ctx.Writer, err.Error(), 500)
	}
}

func (ctx *Context) HTML(code int, html string) {
	ctx.Header("Content-Type", "text/html")
	ctx.Status(code)
	ctx.Writer.Write([]byte(html))
}

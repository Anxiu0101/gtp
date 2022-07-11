package main

import (
	"gtp/gtp"
	"net/http"
)

func main() {
	r := gtp.New()
	r.GET("/", func(ctx *gtp.Context) {
		ctx.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	r.GET("/hello", func(ctx *gtp.Context) {
		// expect /hello?name=anxiu
		ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
	})

	r.GET("/hello/:name", func(ctx *gtp.Context) {
		// expect /hello/anxiu
		ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Param("name"), ctx.Path)
	})

	r.GET("/assets/*filepath", func(ctx *gtp.Context) {
		ctx.JSON(http.StatusOK, gtp.H{"filepath": ctx.Param("filepath")})
	})

	r.Run(":8080")
}

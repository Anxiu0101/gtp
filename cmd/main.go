package main

import (
	"gtp/gtp"
	"net/http"
)

func main() {
	r := gtp.Default()

	r.GET("/", func(c *gtp.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	v1 := r.Group("/v1")
	v1.Use(gtp.Cors()) // global middleware
	{
		v1.GET("/:name", func(ctx *gtp.Context) {
			ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
		})
	}

	v2 := r.Group("/v2")
	{
		v2.GET("/:name", func(ctx *gtp.Context) {
			ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
		})
		v2.POST("/:name", func(ctx *gtp.Context) {
			ctx.JSON(http.StatusOK, gtp.H{
				"username": ctx.PostForm("username"),
				"password": ctx.PostForm("password"),
			})
		})
		v2.DELETE("/:name", func(ctx *gtp.Context) {
			ctx.String(http.StatusOK, "hello %s, delete test OK", ctx.Query("name"))
		})
		v2.PUT("/:name", func(ctx *gtp.Context) {
			ctx.String(http.StatusOK, "hello %s, PUT test OK", ctx.Query("name"))
		})
	}

	// index out of range for testing Recovery()
	r.GET("/panic", func(ctx *gtp.Context) {
		names := []string{"Anxiu"}
		ctx.String(http.StatusOK, names[2])
	})

	r.Run(":8080")
}

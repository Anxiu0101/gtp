package main

import (
	"gtp/gtp"
	"net/http"
)

func main() {
	r := gtp.Default()
	r.GET("/", func(c *gtp.Context) {
		c.String(http.StatusOK, "Hello Geektutu\n")
	})
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *gtp.Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":8080")
}

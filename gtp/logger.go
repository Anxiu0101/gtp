package gtp

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("[gtp] [%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

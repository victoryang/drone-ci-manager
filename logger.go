package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger formatter
func Logger() gin.HandlerFunc {

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.RequestURI

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		username, _ := c.Get("username")
		if username == nil {
			username = "unknown"
		}
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		comment := c.Errors.String()
		fmt.Fprintf(os.Stdout, "%v|%d|%v|%s|%s|%s|%v\n",
			end.Format("2006-01-02 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
			strings.TrimSpace(comment),
		)

	}
}
package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

// Logger is a middleware that logs incoming HTTP requests.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next() // Call the next handler
		log.Printf("Response status: %d", c.Writer.Status())
	}
}
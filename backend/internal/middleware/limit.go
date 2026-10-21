package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DefaultMaxBodyBytes is the request body cap applied to the whole API.
const DefaultMaxBodyBytes int64 = 256 << 10 // 256KB

// MaxBodyBytes limits request bodies to max bytes. Requests that declare a
// larger Content-Length are rejected up front with 413; for undeclared
// lengths (chunked) the body reader fails once the limit is exceeded.
func MaxBodyBytes(max int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > max {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{"error": "Request body too large"})
			return
		}
		if c.Request.Body != nil {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, max)
		}
		c.Next()
	}
}

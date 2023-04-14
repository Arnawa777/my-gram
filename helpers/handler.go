package helpers

import "github.com/gin-gonic/gin"

// Call this when you want to get Content Type
func GetContentType(c *gin.Context) string {
	return c.Request.Header.Get("Content-type")
}

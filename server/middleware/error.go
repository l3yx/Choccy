package middleware

import "github.com/gin-gonic/gin"

func ErrorMiddleware(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(200, gin.H{"err": err})
			c.Abort()
		}
	}()
	c.Next()
}

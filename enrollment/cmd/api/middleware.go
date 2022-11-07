package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.GetHeader("X-API-KEY") == "" {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"error":   "true",
				"message": "Must have X-API-KEY header!",
			})
			c.AbortWithStatus(401)
			return
		}

		if c.GetHeader("X-API-KEY") != os.Getenv("X_API_KEY") {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"error":   "true",
				"message": "Key is mismatched!",
			})
			c.AbortWithStatus(401)
			return
		}
		c.Next()
	}
}
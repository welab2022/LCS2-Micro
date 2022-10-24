package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

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

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AddCorsHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}

}

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {

		log.Printf("Client IP: %s", c.ClientIP())
		log.Printf("Request header: %s", c.Request.Header)

		if c.Request.Method == "OPTIONS" {
			log.Printf("CORS headers")
		}

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

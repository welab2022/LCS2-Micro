package main

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var webPort = "80"

const GROUP_AUTH_API = "/api/auth/"

func (app *Config) startApp() {

	router := gin.New()

	// Apply the middleware to the router (works on groups too)
	// Set up CORS middleware options
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8000", "http://localhost:3000", "http://localhost:8081", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "PATCH", "POST", "PUT", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Origin", "X-API-Key", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length, Link, X-API-Key"},
		MaxAge:           12 * time.Hour,
		AllowCredentials: true,
	}))

	router.GET("/heartbeat", app.HeartBeat)

	// auth with middleware
	authorized := router.Group("/")
	authorized.Use(AuthMiddleWare())
	{
		// map to URL
		authorized.POST("/signin", app.Signin)
		authorized.POST("/upload", app.UpdateAvatar)
		authorized.GET("/avatar/:email", app.GetAvatar)
		authorized.POST("resetpwd", app.ResetPassword)

		authorized.POST("/logout", app.Logout)
		authorized.POST("/refresh", app.Refresh)
		authorized.POST("/changepwd", app.ChangePassword)
		authorized.POST("/adduser", app.AddUser)
		authorized.GET("/listusers", app.ListAllUsers)
		authorized.GET("/user/:email", app.GetUser)

	}

	if os.Getenv("AUTH_PORT") != "" {
		webPort = os.Getenv("AUTH_PORT")
	}
	router.Run(":" + webPort)
}

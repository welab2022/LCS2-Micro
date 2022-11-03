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
	router.Use(cors.Default())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"OPIONS, GET, PUT, POST, DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length, Link"},
		MaxAge:           12 * time.Hour,
		AllowCredentials: true,
	}))

	router.GET("/heartbeat", app.HeartBeat)
	router.POST("/signin", app.Signin)
	router.POST("/upload", app.saveFileHandler)

	// auth with middleware
	authorized := router.Group("/")
	authorized.Use(AuthMiddleWare())
	{
		// map to URL
		authorized.POST("/logout", app.Logout)
		authorized.POST("/refresh", app.Refresh)
		authorized.POST("/changepwd", app.ChangePassword)
		authorized.POST("/adduser", app.AddUser)
		authorized.GET("/listusers", app.ListAllUsers)
	}

	if os.Getenv("AUTH_PORT") != "" {
		webPort = os.Getenv("AUTH_PORT")
	}
	router.Run(":" + webPort)
}

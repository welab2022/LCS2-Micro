package main

import (
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

const webPort = "80"

func (app *Config) startApp() {

	router := gin.New()

	// Set up CORS middleware options
	config := cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          1 * time.Minute,
		Credentials:     false,
		ValidateHeaders: false,
	}

	// map to URL
	router.GET("/heartbeat", app.HeartBeat)
	router.POST("/signin", app.Signin)

	// Apply the middleware to the router (works on groups too)
	router.Use(cors.Middleware(config))

	authorized := router.Group("/")
	authorized.Use(AuthMiddleWare())
	{
		authorized.POST("/logout", app.Logout)
		authorized.POST("/refresh", app.Refresh)
		router.POST("/adduser", app.AddUser)
	}

	router.Run(":" + webPort)
}

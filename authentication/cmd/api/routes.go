package main

import (
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

const webPort = "80"
const SECRET_STRING = "secret_lcs2"
const SESSION_STRING = "session_lcs2"

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
	router.POST("/login", app.LoginHandler)

	// Apply the middleware to the router (works on groups too)
	router.Use(HeartBeatHandler())
	router.Use(cors.Middleware(config))

	authorized := router.Group("/")
	authorized.Use(AuthMiddleWare())
	{
		authorized.POST("/logout", app.LogoutHandler)
	}

	router.Run(":" + webPort)
}

package main

import (
	"net/http"
	"time"

	cors "github.com/itsjamie/gin-cors"

	"github.com/gin-gonic/gin"
)

func (app *Config) routes() http.Handler {

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
	router.POST("/login", h.LoginHandler)

	// Apply the middleware to the router (works on groups too)
	router.Use(cors.Middleware(config))
	authorized := router.Group("/")
	authorized.Use(AuthMiddleWare())
	{
		authorized.GET("/logout", app.LogoutHandler)
	}

	return router
}

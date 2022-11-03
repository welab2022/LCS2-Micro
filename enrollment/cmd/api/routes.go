package main

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

const GROUP_ENROL_API = "/api/enroll/"

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

	// auth with middleware
	authorized := router.Group("/")
	authorized.Use(AuthMiddleWare())
	{
		// map to URL
		authorized.GET("/list", app.ListEnroll)
	}

	if os.Getenv("ENROLL_PORT") != "" {
		webPort = os.Getenv("ENROLL_PORT")
	}
	router.Run(":" + webPort)
}

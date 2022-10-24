package main

import (
	"github.com/gin-gonic/gin"
)

const webPort = "80"

func (app *Config) startApp() {

	router := gin.New()

	// Set up CORS middleware options
	// config := cors.Config{
	// 	Origins:         "*",
	// 	Methods:         "OPIONS, GET, PUT, POST, DELETE",
	// 	RequestHeaders:  "Origin, Authorization, Content-Type, Accept-Encoding, X-CSRF-Token",
	// 	ExposedHeaders:  "",
	// 	MaxAge:          12 * time.Hour,
	// 	Credentials:     false,
	// 	ValidateHeaders: false,
	// }

	// map to URL
	router.GET("/heartbeat", app.HeartBeat)
	router.POST("/signin", app.Signin)

	// Apply the middleware to the router (works on groups too)
	// router.Use(cors.Middleware(config))
	router.Use(CORSMiddleware())

	authorized := router.Group("/")
	authorized.Use(AuthMiddleWare())
	{
		authorized.POST("/logout", app.Logout)
		authorized.POST("/refresh", app.Refresh)
		router.POST("/adduser", app.AddUser)
	}

	router.Run(":" + webPort)
}

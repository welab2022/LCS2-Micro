package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/welab2022/LCS2-Micro/authentication/handlers"
)

const webPort = "80"

var (
	router = gin.New()
)

func main() {
	log.Println("Authentication service start...")

	router.GET("/login", handlers.LoginHandler)
	router.Run(":" + webPort)
}

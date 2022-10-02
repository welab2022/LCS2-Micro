package main

import (
	"log"

	"heartbeat/handlers"

	"github.com/gin-gonic/gin"
)

const webPort = "80"

var (
	router = gin.Default()
)

func main() {
	log.Println("Heartbeat service start...")

	router.GET("/heartbeat", handlers.HeartBeatHandler)
	router.Run(":" + webPort)
}

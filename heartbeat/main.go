package main

import (
	"log"

	"github.com/welab2022/LCS2-Micro/heartbeat/handlers"
	"github.com/gin-gonic/gin"
)

const webPort = "80"

var (
	router = gin.New()
)

func main() {
	log.Println("Heartbeat service start...")

	router.GET("/heartbeat", handlers.HeartBeatHandler)
	router.Run(":" + webPort)
}

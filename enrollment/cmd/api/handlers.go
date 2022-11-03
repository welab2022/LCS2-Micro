package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *Config) HeartBeat(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "200",
		"title":   "Health OK",
		"updated": GetDateString(),
	})
}

func (app *Config) ListEnroll(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "200",
		"message": "List enrolls",
	})
}

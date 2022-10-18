package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type jsonResponse struct {
	Status  string
	Message string
	Token   string
	Data    interface{}
}

// swagger:operation GET /login to login the system
// Check authentication
// ---
// produces:
// - application/json
// parameters:
//
// responses:
//
//	'200':
//	    description: Successful operation
//	'503':
//	    description: Service not found
func (app *Config) LoginHandler(ctx *gin.Context) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	ctx.Header("Content-Type", "application/json; charset=utf-8")

	if err := ctx.ShouldBindJSON(&requestPayload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validate the user against database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Invalid email"})
		return
	}

	if requestPayload.Password != "letmein" {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Invalid password"})
		return
	}

	token := GenerateTokenBase64()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	payload := jsonResponse{
		Status:  "success",
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Token:   token,
		Data:    user,
	}

	ctx.JSON(http.StatusAccepted, payload)
}

// swagger:operation GET /login to login the system
// Check authentication
// ---
// produces:
// - application/json
// parameters:
//
// responses:
//
//	'200':
//	    description: Successful operation
//	'503':
//	    description: Service not found
func (app *Config) LogoutHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "logout",
	})

}

// swagger:operation GET /heartbeat get heartbeat
// Get heartbeat
// ---
// produces:
// - application/json
// parameters:
//
// responses:
//
//	'200':
//	    description: Successful operation
//	'503':
//	    description: Service not found
func (app *Config) HeartBeatHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{
		"status":  "200",
		"title":   "Health OK",
		"updated": GetDateString(),
	})
}

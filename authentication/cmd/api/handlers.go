package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type jsonResponse struct {
	Status  string
	Message string
	Token   string
	Data    interface{}
}

// each session contains the username of the user and the time at which it expires
type session struct {
	username string
	expiry   time.Time
}

var sessions = map[string]session{}

// we'll use this method later to determine if the session has expired
func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

// Create a struct that models the structure of a user in the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
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

	var creds Credentials

	ctx.Header("Content-Type", "application/json; charset=utf-8")

	if err := ctx.ShouldBindJSON(&requestPayload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "true",
			"message": err.Error(),
		})
		return
	}

	// validate the user against database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Invalid email"})
		return
	}

	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session token"})
	// 	return
	// }

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Invalid password"})
		return
	}

	log.Printf("Credential is correct!")
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	// Set the token in the session map, along with the user whom it represents
	sessions[sessionToken] = session{
		username: creds.Username,
		expiry:   expiresAt,
	}

	token := GenerateTokenBase64()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	payload := jsonResponse{
		Status:  "success",
		Message: fmt.Sprintf("Authenticated! Logged in user: %s", user.Email),
		Token:   token,
		Data:    user,
	}

	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 120 seconds
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})

	log.Printf("session_token: %s", sessionToken)

	ctx.JSON(http.StatusAccepted, payload)
}

// swagger:operation GET /logout to login the system
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
	c, err := ctx.Request.Cookie("session_token")

	log.Printf("cookie: %s", c.Value)

	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		// For any other type of error, return a bad request status
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	sessionToken := c.Value

	// remove the users session from the session map
	delete(sessions, sessionToken)

	// We need to let the client know that the cookie is expired
	// In the response, we set the session token to an empty
	// value and set its expiry as the current time
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out!"})
}

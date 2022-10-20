package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type jsonResponse struct {
	Status  string
	Message string
	API     string
	Data    interface{}
}

// each session contains the username of the user and the time at which it expires
type session struct {
	username string
	expiry   time.Time
}

var sessions = map[string]session{}

const SESSION_TOKEN = "lcs2_session_token"

// we'll use this method later to determine if the session has expired
func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

// Create a struct that models the structure of a user in the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func (app *Config) HeartBeat(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "200",
		"title":   "Health OK",
		"updated": GetDateString(),
	})
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
func (app *Config) Signin(ctx *gin.Context) {

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

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Invalid password"})
		return
	}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	// Set the token in the session map, along with the user whom it represents
	sessions[sessionToken] = session{
		username: creds.Username,
		expiry:   expiresAt,
	}

	api_key := GenerateTokenBase64()

	payload := jsonResponse{
		Status:  "success",
		Message: fmt.Sprintf("Authenticated! Logged in user: %s", user.Email),
		API:     api_key,
		Data:    user,
	}

	// Finally, we set the client cookie for SESSION_TOKEN as the session token we just generated
	// we also set an expiry time of 120 seconds
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    SESSION_TOKEN,
		Value:   sessionToken,
		Expires: expiresAt,
	})

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
func (app *Config) Logout(ctx *gin.Context) {

	c, err := ctx.Request.Cookie(SESSION_TOKEN)
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

	// We then get the name of the user from our session map, where we set the session token
	userSession, exists := sessions[sessionToken]
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if userSession.isExpired() {
		delete(sessions, sessionToken)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// remove the users session from the session map
	delete(sessions, sessionToken)

	// We need to let the client know that the cookie is expired
	// In the response, we set the session token to an empty
	// value and set its expiry as the current time
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    SESSION_TOKEN,
		Value:   "",
		Expires: time.Now(),
	})

	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out!"})
}

func (app *Config) Register(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Register ok!"})
}

func (app *Config) Refresh(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "session refreshed ok!"})
}

package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/welab2022/LCS2-Micro/authentication/data"
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
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Expired => Unauthorized. Please refresh your session!"})
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

func (app *Config) validateSession(ctx *gin.Context) bool {
	c, err := ctx.Request.Cookie(SESSION_TOKEN)
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return false
		}
		// For any other type of error, return a bad request status
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return false
	}

	sessionToken := c.Value

	// We then get the name of the user from our session map, where we set the session token
	userSession, exists := sessions[sessionToken]
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return false
	}

	if userSession.isExpired() {
		delete(sessions, sessionToken)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Expired => Unauthorized. Please refresh your session!"})
		return false
	}
	return true
}

func (app *Config) AddUser(ctx *gin.Context) {

	if !app.validateSession(ctx) {
		return
	}

	// user info request
	var requestPayload struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name,omitempty"`
		LastName  string `json:"last_name,omitempty"`
		Password  string `json:"password"`
	}

	var responseUser struct {
		Status  string
		Message string
	}

	ctx.Header("Content-Type", "application/json; charset=utf-8")

	if err := ctx.ShouldBindJSON(&requestPayload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "true",
			"message": err.Error(),
		})
		return
	}

	// check if the user against database is existed
	_, err := app.Models.User.GetByEmail(requestPayload.Email)

	if err == nil {
		responseUser.Status = "error"
		responseUser.Message = fmt.Sprintf("User %s existed!", requestPayload.Email)
		ctx.JSON(http.StatusNotFound, responseUser)
		return
	}

	var user data.User
	user.Email = requestPayload.Email
	user.FirstName = requestPayload.FirstName
	user.LastName = requestPayload.LastName
	user.Password = requestPayload.Password
	user.Active = 1
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	id, err := app.Models.User.Insert(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal error, db access failed!"})
		return
	}

	responseUser.Status = "User added!"
	responseUser.Message = fmt.Sprintf("User %s added and id: %d!", requestPayload.Email, id)
	ctx.JSON(http.StatusOK, responseUser)

}

func (app *Config) Refresh(ctx *gin.Context) {

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

	// If the previous session is valid, create a new session token for the current user
	newSessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	// Set the token in the session map, along with the user whom it represents
	sessions[newSessionToken] = session{
		username: userSession.username,
		expiry:   expiresAt,
	}

	// Delete the older session token
	delete(sessions, sessionToken)

	// Finally, we set the client cookie for SESSION_TOKEN as the session token we just generated
	// we also set an expiry time of 120 seconds
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    SESSION_TOKEN,
		Value:   newSessionToken,
		Expires: expiresAt,
	})

	ctx.JSON(http.StatusOK, gin.H{"message": "Session refreshed OK!"})
}

func (app *Config) ChangePassword(ctx *gin.Context) {
	if !app.validateSession(ctx) {
		return
	}

	// user info request for changing password
	var requestPayload struct {
		Email       string `json:"email"`
		OldPassword string `json:"old_password,omitempty"`
		NewPassword string `json:"new_password,omitempty"`
	}

	var responseUser struct {
		Status  string
		Message string
	}

	ctx.Header("Content-Type", "application/json; charset=utf-8")

	if err := ctx.ShouldBindJSON(&requestPayload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "true",
			"message": err.Error(),
		})
		return
	}

	// check if the user against database is existed
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		// the user doesn't exists
		responseUser.Status = "error"
		responseUser.Message = fmt.Sprintf("User %s doesn't exist!", requestPayload.Email)
		ctx.JSON(http.StatusNotFound, responseUser)
		return
	}

	// the user exists and then check if the old password is matched
	valid, err := user.PasswordMatches(requestPayload.OldPassword)
	if err != nil || !valid {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Invalid password"})
		return
	}

	err = user.ResetPassword(requestPayload.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Change password failure"})
		return
	}

	responseUser.Status = "OK"
	responseUser.Message = "Password changed successfully"
	ctx.JSON(http.StatusOK, responseUser)
}

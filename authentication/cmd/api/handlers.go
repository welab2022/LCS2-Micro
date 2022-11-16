package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/welab2022/LCS2-Micro/authentication/data"
)

type mailMessage struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

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
	expiresAt := time.Now().Add(1000 * 60 * 60 * 24 * 30 * time.Second) // 30 days

	// Set the token in the session map, along with the user whom it represents
	sessions[sessionToken] = session{
		username: requestPayload.Email,
		expiry:   expiresAt,
	}

	api_key := GenerateTokenBase64()

	payload := jsonResponse{
		Status:  "success",
		Message: fmt.Sprintf("Authenticated! Logged in user: %s", user.Email),
		API:     api_key,
		Data:    user,
	}

	// update the last_login
	now := time.Now()
	err = user.LastLoginUpdate(now, user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Can't not update database"})
		return
	}

	// Finally, we set the client cookie for SESSION_TOKEN as the session token we just generated
	// we also set an expiry time of 120 seconds
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     SESSION_TOKEN,
		Value:    sessionToken,
		Expires:  expiresAt,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
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

	sessionToken, err := app.validateSession(ctx)
	if err != nil {
		return
	}

	var requestPayload struct {
		Email string `json:"email"`
	}

	ctx.Header("Content-Type", "application/json; charset=utf-8")

	if err := ctx.ShouldBindJSON(&requestPayload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "true",
			"message": err.Error(),
		})
		return
	}

	if requestPayload.Email != sessions[sessionToken].username {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "true",
			"message": "You must to sign in before log out",
		})
		return
	}

	// remove the users session from the session map
	delete(sessions, sessionToken)

	// We need to let the client know that the cookie is expired
	// In the response, we set the session token to an empty
	// value and set its expiry as the current time
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     SESSION_TOKEN,
		Value:    "",
		Expires:  time.Now(),
		SameSite: http.SameSiteDefaultMode,
	})

	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out!"})
}

func (app *Config) validateSession(ctx *gin.Context) (string, error) {

	c, err := ctx.Request.Cookie(SESSION_TOKEN)
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return "", err
		}
		// For any other type of error, return a bad request status
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return "", err
	}

	sessionToken := c.Value

	// We then get the name of the user from our session map, where we set the session token
	userSession, exists := sessions[sessionToken]
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return "", err
	}

	if userSession.isExpired() {
		delete(sessions, sessionToken)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Expired => Unauthorized. Please refresh your session!"})
		return "", err
	}

	return sessionToken, nil
}

func (app *Config) AddUser(ctx *gin.Context) {

	sessionToken, err := app.validateSession(ctx)
	if err != nil {
		return
	}

	if sessions[sessionToken].username != "admin@example.com" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "No permission to add a new user"})
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

	if err = ctx.ShouldBindJSON(&requestPayload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "true",
			"message": err.Error(),
		})
		log.Printf("bindJsonReq error: %s", err)
		return
	}

	// check if the user against database is existed
	exist_user, err := app.Models.User.GetByEmail(requestPayload.Email)

	if exist_user != nil && err == nil {
		responseUser.Status = "error"
		responseUser.Message = fmt.Sprintf("User %s already existed!", requestPayload.Email)
		ctx.JSON(http.StatusInternalServerError, responseUser)
		return
	}
	// password needs to be generated and sent through a user registration email
	var user data.User
	user.Email = requestPayload.Email
	user.FirstName = requestPayload.FirstName
	user.LastName = requestPayload.LastName
	user.Password = requestPayload.Password
	user.Active = 1
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.LastLogin = time.Time{}
	user.PasswordChangeAt = time.Time{}

	id, err := app.Models.User.Insert(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal error, db access failed!"})
		return
	}

	responseUser.Status = "User added!"
	responseUser.Message = fmt.Sprintf("User %s added and id: %d!", requestPayload.Email, id)

	// need to send an email notification
	// ...

	ctx.JSON(http.StatusOK, responseUser)
}

func (app *Config) Refresh(ctx *gin.Context) {

	sessionToken, err := app.validateSession(ctx)
	if err != nil {
		return
	}

	// We then get the name of the user from our session map, where we set the session token
	userSession, exists := sessions[sessionToken]
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// If the previous session is valid, create a new session token for the current user
	newSessionToken := uuid.NewString()
	expiresAt := time.Now().Add(1000 * 60 * 60 * 24 * 30 * time.Second)

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
		Name:     SESSION_TOKEN,
		Value:    newSessionToken,
		Expires:  expiresAt,
		SameSite: http.SameSiteNoneMode,
	})

	ctx.JSON(http.StatusOK, gin.H{"message": "Session refreshed OK!"})
}

func (app *Config) ChangePassword(ctx *gin.Context) {
	sessionToken, err := app.validateSession(ctx)
	if err != nil {
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

	if sessions[sessionToken].username != requestPayload.Email {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "%s is not authenticated yet, please sign in to the system"})
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

func (app *Config) UpdateAvatar(ctx *gin.Context) {
	sessionToken, err := app.validateSession(ctx)
	if err != nil {
		return
	}

	ctx.Header("Content-Type", "application/json; charset=utf-8")

	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		log.Printf("No file is received err: %s", err)
		return
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return
	}

	// need to be authorized here, session to get the email
	email := sessions[sessionToken].username

	err = app.Models.User.UpdateAvatar(buf.Bytes(), email)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to store database",
		})
		return
	}

	// File saved successfully. Return proper result
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Your file has been successfully uploaded.",
	})
}

func (app *Config) GetAvatar(ctx *gin.Context) {
	sessionToken, err := app.validateSession(ctx)
	if err != nil {
		return
	}

	ctx.Header("Content-Type", "application/json; charset=utf-8")

	email := ctx.Param("email")
	if sessions[sessionToken].username != email {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "No permission, unthorized",
		})
		return

	}

	bytes, err := app.Models.User.GetAvatar(email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Failure!",
			"message": fmt.Sprintf("Query %s avatar failed!", email),
		})
		return
	}

	var responseAvatar struct {
		Email    string
		MimeType string
		Base64   string
	}

	var base64Encoding string
	mimeType := http.DetectContentType(bytes)

	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}

	base64Encoding += ToBase64(bytes)

	responseAvatar.Email = email
	responseAvatar.MimeType = mimeType
	responseAvatar.Base64 = base64Encoding

	ctx.JSON(http.StatusOK, responseAvatar)
}

func (app *Config) ListAllUsers(ctx *gin.Context) {
	_, err := app.validateSession(ctx)
	if err != nil {
		return
	}

	ctx.Header("Content-Type", "application/json; charset=utf-8")

	users, err := app.Models.User.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal error, db access failed!"})
		return
	}

	ctx.JSON(http.StatusOK, users)

}

func (app *Config) ResetPassword(ctx *gin.Context) {
	_, err := app.validateSession(ctx)
	if err != nil {
		return
	}

	// user info request for changing password
	var requestPayload struct {
		Email string `json:"email"`
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

	// generate the reset password
	rstpassword, _ := GeneratePassword()

	log.Printf("reset password: %s", rstpassword)

	err = user.ResetPassword(rstpassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Can't reset password"})
		return
	}

	responseUser.Status = "Reset succeeded"
	responseUser.Message = fmt.Sprintf("Reset password is sent to email: %s", requestPayload.Email)

	// sent email
	var mail mailMessage
	mail.From = "admin@example.com"
	mail.To = requestPayload.Email
	mail.Subject = "Reset password"

	mail.Message = fmt.Sprintf("Hello,\n reset password is %s", rstpassword)

	err = app.sendMail(mail)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Can't send reset password email"})
		return
	}

	ctx.JSON(http.StatusOK, responseUser)

}

func (app *Config) sendMail(mail mailMessage) error {

	mailServiceURL := "http://host.docker.internal:9001/send"
	jsonData, _ := json.MarshalIndent(mail, "", "\t")
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("newMessage: sendMail failed %s", err)
		return err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil || response.StatusCode != 202 {
		log.Printf("http: sendMail failed %s", err)
		return err
	}

	return nil
}

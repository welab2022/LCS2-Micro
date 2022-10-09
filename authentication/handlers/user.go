package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	apiDateLayout = "2006-01-02T15:04:05Z"
	apiDBLayout   = "2006-01-02 15:04:05"
)

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
func LoginHandler(ctx *gin.Context) {
	fmt.Println("login")
	ctx.JSON(http.StatusOK, map[string]string{
		"status":  "authenticated",
		"title":   "Login",
		"updated": GetDateString(),
	})
}

func GetNow() time.Time {
	return time.Now().UTC()
}

func GetDateString() string {
	return GetNow().Format(apiDateLayout)
}

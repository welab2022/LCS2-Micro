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
func HeartBeatHandler(ctx *gin.Context) {
	fmt.Println("Heartbeat")
	ctx.JSON(http.StatusOK, map[string]string{
		"status":  "200",
		"title":   "Health OK",
		"updated": GetDateString(),
	})
}

func GetNow() time.Time {
	return time.Now().UTC()
}

func GetDateString() string {
	return GetNow().Format(apiDateLayout)
}

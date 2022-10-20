package main

import (
	"os"
	"time"
)

const (
	apiDateLayout = "2006-01-02T15:04:05Z"
	apiDBLayout   = "2006-01-02 15:04:05"
)

func GetNow() time.Time {
	return time.Now().UTC()
}

func GetDateString() string {
	return GetNow().Format(apiDateLayout)
}

func GetNowDBDate() string {
	return GetNow().Format(apiDBLayout)
}

func GenerateTokenBase64() string {

	return os.Getenv("X_API_KEY")
}

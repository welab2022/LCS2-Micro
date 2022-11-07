package main

import (
	"encoding/base64"
	"os"
	"time"

	_ "github.com/nicored/avatar"
	"github.com/sethvargo/go-password/password"
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

func ToBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func GeneratePassword() (string, error) {
	// This is exactly the same as calling "Generate" directly. It will use all
	// the default values.
	// gen, err := password.NewGenerator(nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// _ = gen // gen.Generate(...)
	// gen.Generate()
	var (
		length      = 6
		numDigits   = 2
		numSymbols  = 2
		noUpper     = false
		allowRepeat = true
	)
	return password.Generate(length, numDigits, numSymbols, noUpper, allowRepeat)
}

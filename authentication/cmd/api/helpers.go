package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"log"
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

func ReadBinFile(binfile string) (string, error) {
	f, err := os.Open(binfile)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	reader := bufio.NewReader(f)
	buf := make([]byte, 256)
	output := make([]byte, 1024)

	for {
		_, err := reader.Read(buf)

		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}

		fmt.Printf("%s", hex.Dump(buf))
		output = append(output[:], buf[:]...)
	}
	fmt.Println("got file:")
	fmt.Printf("%s", hex.Dump(output))
	return hex.Dump(output), nil
}

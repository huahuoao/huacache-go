package client

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"mime/multipart"
	"net/http"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(length int) string {
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
		if err != nil {
			return ""
		}
		result[i] = letterBytes[num.Int64()]
	}
	return string(result)
}
func HttpTest() {
	// Prepare the data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	writer.WriteField("key", generateRandomString(10))
	writer.WriteField("value", generateRandomString(10))
	writer.WriteField("group", "huahuo1")
	writer.Close()

	// Create the request
	req, err := http.NewRequest("POST", "http://localhost:4160/huacache/set", &requestBody)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
}

func FiveThousandTest() {
	start := time.Now()
	for i := 0; i < 5000; i++ {
		HttpTest()
	}
	duration := time.Since(start).Milliseconds() // 获取毫秒
	fmt.Printf("HTTP FiveThousandTest done: %.2f ms\n", float64(duration))
}

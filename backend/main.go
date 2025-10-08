package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type TokenRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
}

func waitForAPI(apiURL string, maxRetries int) error {
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get(apiURL + "/health")
		if err == nil && resp.StatusCode == 200 {
			log.Println("API is ready!")
			return nil
		}
		log.Printf("Waiting for API... (attempt %d/%d)", i+1, maxRetries)
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("API not ready after %d attempts", maxRetries)
}

func main() {

	// Warte bis API bereit ist
	if err := waitForAPI("http://api:8080", 10); err != nil {
		log.Fatal(err)
	}

	/*
		filePath := flag.String("filepath", "", "Path to the file")

		flag.Parse()

		fmt.Println("File path:", *filePath)

		if *filePath == "" {
			fmt.Println("No filepath provided")
			return
		} */

	filePath := "/home/marcel/dev/scripts/go/backend/video.mp4"

	originalRefreshToken := "rft.7yekSfYUqyhHt7f6Inz3wkJ9ErZZ0lZkbuFrejf5n0KuKYXZcL13x3GqTuZV!4736.e1"

	// Load German timezone (Europe/Berlin)
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}

	// Renew access token and get expiry in seconds
	accessToken, expiresIn, err := RenewAccessToken(originalRefreshToken)
	if err != nil {
		panic(err)
	}

	// Calculate expires_at using TikTok's expires_in
	expiresAt := time.Now().In(loc).Add(time.Duration(expiresIn) * time.Second).Format(time.RFC3339)

	// Build the request payload
	payload := TokenRequest{
		AccessToken:  accessToken,
		RefreshToken: originalRefreshToken, // optionally replace with new refresh token if returned
		ExpiresAt:    expiresAt,
	}

	// Encode to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	/*
			// Send POST request to API
		resp, err := http.Post(
			"http://127.0.0.1:8080/api/token/create",
			"application/json",
			bytes.NewBuffer(jsonData),
		)

	*/

	// Send POST request to API
	resp, err := http.Post(
		"http://api:8080/api/token/create",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	// Create Upload URL for File

	contentType := "video/mp2"
	uploadUrl, err := CreateUploadURL(
		"Test",
		"SELF_ONLY",
		filePath,
		1000,
		accessToken,
		originalRefreshToken,
	)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("URL:", uploadUrl)

	fileSize, _, err := GetFileSize(filePath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	err = UploadFileComplete(uploadUrl, filePath, fileSize, contentType)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

}

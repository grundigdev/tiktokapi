package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TokenRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
}

func main() {
	originalRefreshToken := "rft.7yekSfYUqyhHt7f6Inz3wkJ9ErZZ0lZkbuFrejf5n0KuKYXZcL13x3GqTuZV!4736.e1"

	// Renew access token and get expiry in seconds
	accessToken, expiresIn, err := RenewAccessToken(originalRefreshToken)
	if err != nil {
		panic(err)
	}

	fmt.Println("New Access Token:", accessToken)

	// Load German timezone (Europe/Berlin)
	loc, err := time.LoadLocation("Europe/Berlin")
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

	// Send POST request to your API
	resp, err := http.Post(
		"http://127.0.0.1:8080/api/token/create",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	// Create Upload URL for File

	uploadUrl, err := CreateUploadURL(
		"Test",
		"SELF_ONLY",
		1000,
		1803949,
		1803949,
		1,
		accessToken,
		originalRefreshToken,
	)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("URL:", uploadUrl)
}

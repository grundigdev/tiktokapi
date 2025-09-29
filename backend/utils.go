package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Response struct based on TikTok's OAuth response format
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"` // changed to int64 for time.Duration
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func RenewAccessToken(refreshToken string) (string, int64, error) {
	endpoint := "https://open.tiktokapis.com/v2/oauth/token/"

	// Prepare form data
	data := url.Values{}
	data.Set("client_key", "sbawwttwjocxh640nv")
	data.Set("client_secret", "xTzcl06NLOiEmSyGBSr5TlDAQCVUsx4a")
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	// Create request
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return "", 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cache-Control", "no-cache")

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("failed to read response: %w", err)
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", 0, fmt.Errorf("failed to parse response: %w\nResponse body: %s", err, string(body))
	}

	if tokenResp.AccessToken == "" {
		return "", 0, fmt.Errorf("no access_token in response: %s", string(body))
	}

	return tokenResp.AccessToken, tokenResp.ExpiresIn, nil
}

// CHECK IF ACCESS TOKEN IS EXPIRED OR VALID

type CheckTokenRequest struct {
	AccessToken string `json:"access_token"`
}

type TokenAPIResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Data    TokenData `json:"data"`
}

type TokenData struct {
	Token   interface{} `json:"token"`
	IsValid bool        `json:"isValid"`
}

// Fetches the latest token and renews it if expired
func CheckAccessToken(accessToken string) (bool, error) {
	endpoint := "http://127.0.0.1:8080/api/token/get"

	body := CheckTokenRequest{
		AccessToken: accessToken,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return false, fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("received non-OK response: %s", resp.Status)
	}

	// Decode the response body
	var checkTokenResp TokenAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&checkTokenResp); err != nil {
		return false, fmt.Errorf("failed to decode response: %v", err)
	}

	// Check if isValid is true
	if !checkTokenResp.Data.IsValid {
		return false, fmt.Errorf("token is not valid")
	}

	return true, nil
}

// CREATE UPLOAD URL

type CreateUploadURLPostInfo struct {
	Title                 string `json:"title"`
	PrivacyLevel          string `json:"privacy_level"`
	DisableDuet           bool   `json:"disable_duet"`
	DisableComment        bool   `json:"disable_comment"`
	DisableStitch         bool   `json:"disable_stitch"`
	VideoCoverTimestampMS int    `json:"video_cover_timestamp_ms"`
}

type CreateUploadURLSourceInfo struct {
	Source          string `json:"source"`
	VideoSize       int    `json:"video_size"`
	ChunkSize       int    `json:"chunk_size"`
	TotalChunkCount int    `json:"total_chunk_count"`
}

type CreateUploadURLRequestBody struct {
	PostInfo   CreateUploadURLPostInfo   `json:"post_info"`
	SourceInfo CreateUploadURLSourceInfo `json:"source_info"`
}

type CreateUploadURLResponse struct {
	Data struct {
		PublishID string `json:"publish_id"`
		UploadURL string `json:"upload_url"`
	} `json:"data"`
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		LogID   string `json:"log_id"`
	} `json:"error"`
}

func CreateUploadURL(
	title string,
	privacyLevel string,
	videoCoverTimestampMS int,
	videoSize int,
	chunkSize int,
	totalChunkCount int,
	accessToken string,
	refreshToken string,
) (string, error) {

	// Check if access token is valid
	isValid, err := CheckAccessToken(accessToken)
	if err != nil {
		return "", fmt.Errorf("failed to check access token: %v", err)
	}

	// Renew token if invalid
	if !isValid {
		newAccessToken, _, err := RenewAccessToken(refreshToken)
		if err != nil {
			return "", fmt.Errorf("failed to renew access token: %v", err)
		}
		accessToken = newAccessToken
	}

	url := "https://open.tiktokapis.com/v2/post/publish/video/init/"

	body := CreateUploadURLRequestBody{
		PostInfo: CreateUploadURLPostInfo{
			Title:                 title,
			PrivacyLevel:          privacyLevel,
			DisableDuet:           true,
			DisableComment:        true,
			DisableStitch:         true,
			VideoCoverTimestampMS: videoCoverTimestampMS,
		},
		SourceInfo: CreateUploadURLSourceInfo{
			Source:          "FILE_UPLOAD",
			VideoSize:       videoSize,
			ChunkSize:       chunkSize,
			TotalChunkCount: totalChunkCount,
		},
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-OK response: %s", resp.Status)
	}

	// Response Body lesen
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-OK response: %s, body: %s", resp.Status, string(bodyBytes))
	}

	// Parse JSON Response
	var uploadResponse CreateUploadURLResponse
	if err := json.Unmarshal(bodyBytes, &uploadResponse); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	// Check if error occurred
	if uploadResponse.Error.Code != "ok" {
		return "", fmt.Errorf("API error: %s - %s", uploadResponse.Error.Code, uploadResponse.Error.Message)
	}

	fmt.Println("Video publish init request successful")

	return uploadResponse.Data.UploadURL, nil

}

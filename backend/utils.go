package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
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

// GET FILE SIZE

func GetFileSize(filepath string) (int64, string, error) {
	// Check if file exists and get file info
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, "", fmt.Errorf("file does not exist: %s", filepath)
		}
		return 0, "", fmt.Errorf("failed to access file %s: %w", filepath, err)
	}

	// Check if path is a directory
	if fileInfo.IsDir() {
		return 0, "", fmt.Errorf("path is a directory, not a file: %s", filepath)
	}

	// Detect file type
	file, err := os.Open(filepath)
	if err != nil {
		return 0, "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read first 512 bytes for type detection
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil {
		return 0, "", fmt.Errorf("failed to read file: %w", err)
	}

	// Detect content type
	contentType := http.DetectContentType(buffer[:n])

	// Validate against allowed video formats
	allowedFormats := map[string]bool{
		"video/mp4":       true,
		"video/quicktime": true,
		"video/webm":      true,
	}

	if !allowedFormats[contentType] {
		return 0, "", fmt.Errorf("invalid file format: %s (expected video/mp4, video/quicktime, or video/webm)", contentType)
	}

	return fileInfo.Size(), contentType, nil
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
	VideoSize       int64  `json:"video_size"`
	ChunkSize       int64  `json:"chunk_size"`
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
	filePath string,
	videoCoverTimestampMS int,
	accessToken string,
	refreshToken string,
) (string, error) {

	var singleUpload bool
	fileSize, _, err := GetFileSize(filePath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return "", fmt.Errorf("failed to get file size: %v", err)
	}

	if fileSize > 5000000 {
		singleUpload = false
	} else {
		singleUpload = true
	}

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

	if singleUpload {
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
				VideoSize:       fileSize,
				ChunkSize:       fileSize,
				TotalChunkCount: 1,
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
	} else {
		// Calculate chunk size and count for multi-chunk upload
		// Chunk size must be between 5 MB and 64 MB
		// const minChunkSize = 5000000  // 5 MB
		const maxChunkSize = 64000000 // 64 MB

		/*
			chunkSize := maxChunkSize // Use maximum chunk size for efficiency


			// If there's a remainder, we need one more chunk
			if fileSize%int64(chunkSize) > 0 {
				totalChunkCount++
			}

			// Ensure we don't exceed 1000 chunks
			if totalChunkCount > 1000 {
				// Recalculate chunk size to fit within 1000 chunks
				chunkSize = int(fileSize / 1000)
				// Round up to ensure we don't exceed 1000 chunks
				if fileSize%1000 > 0 {
					chunkSize++
				}
				totalChunkCount = 1000
			}
		*/

		chunkSize := 10000000
		totalChunkCount := int(fileSize / int64(chunkSize))

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
				VideoSize:       fileSize,
				ChunkSize:       10000000,
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

		fmt.Println("Video publish init request successful (chunked upload)")

		return uploadResponse.Data.UploadURL, nil
	}
}

// UploadFileComplete handles both single and multi-chunk uploads
func UploadFileComplete(
	uploadUrl string,
	filePath string,
	fileSize int64,
	contentType string,
) error {
	const minChunkSize = 5000000  // 5 MB
	const maxChunkSize = 64000000 // 64 MB

	// Determine if single or multi-chunk upload
	if fileSize < minChunkSize {
		return uploadSingleChunk(uploadUrl, filePath, fileSize, contentType)
	} else {
		return uploadMultiChunk(uploadUrl, filePath, fileSize, contentType, maxChunkSize)
	}
}

// uploadSingleChunk handles files < 5 MB
func uploadSingleChunk(uploadUrl string, filePath string, fileSize int64, contentType string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read entire file
	fileData, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("PUT", uploadUrl, bytes.NewReader(fileData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers for single chunk upload
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Length", strconv.FormatInt(fileSize, 10))
	req.Header.Set("Content-Range", fmt.Sprintf("bytes 0-%d/%d", fileSize-1, fileSize))

	// Execute request
	client := &http.Client{
		Timeout: 5 * time.Minute,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Single chunk should return 201
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	fmt.Println("Single chunk upload successful")
	return nil
}

// uploadMultiChunk handles files >= 5 MB
func uploadMultiChunk(uploadUrl string, filePath string, fileSize int64, contentType string, chunkSize int) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Calculate number of chunks
	/*totalChunks := int(fileSize / int64(chunkSize))
	if fileSize%int64(chunkSize) > 0 {
		totalChunks++
	}*/

	chunkSize2 := 10000000
	totalChunks := int(fileSize / int64(chunkSize2))

	fmt.Printf("Starting multi-chunk upload: %d chunks\n", totalChunks)

	// Upload each chunk sequentially
	for chunkIndex := 0; chunkIndex < totalChunks; chunkIndex++ {
		firstByte := int64(chunkIndex) * int64(chunkSize2)
		lastByte := firstByte + int64(chunkSize2) - 1

		// Last chunk: adjust to file size and can exceed chunkSize
		if chunkIndex == totalChunks-1 {
			lastByte = fileSize - 1
		}

		currentChunkSize := lastByte - firstByte + 1

		// Upload this chunk
		err := uploadChunk(uploadUrl, file, firstByte, lastByte, fileSize, currentChunkSize, contentType, chunkIndex+1, totalChunks)
		if err != nil {
			return fmt.Errorf("failed to upload chunk %d/%d: %w", chunkIndex+1, totalChunks, err)
		}

		fmt.Printf("Uploaded chunk %d/%d\n", chunkIndex+1, totalChunks)
	}

	fmt.Println("Multi-chunk upload completed successfully")
	return nil
}

// uploadChunk uploads a single chunk
func uploadChunk(
	uploadUrl string,
	file *os.File,
	firstByte int64,
	lastByte int64,
	totalBytes int64,
	chunkSize int64,
	contentType string,
	chunkNum int,
	totalChunks int,
) error {
	// Seek to starting position
	_, err := file.Seek(firstByte, 0)
	if err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}

	// Read chunk data
	chunkData := make([]byte, chunkSize)
	bytesRead, err := io.ReadFull(file, chunkData)
	if err != nil && err != io.ErrUnexpectedEOF {
		return fmt.Errorf("failed to read chunk: %w", err)
	}

	// Verify we read the expected amount
	if int64(bytesRead) != chunkSize {
		return fmt.Errorf("read %d bytes but expected %d", bytesRead, chunkSize)
	}

	// Create HTTP request
	req, err := http.NewRequest("PUT", uploadUrl, bytes.NewReader(chunkData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	contentRange := fmt.Sprintf("bytes %d-%d/%d", firstByte, lastByte, totalBytes)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Length", strconv.FormatInt(chunkSize, 10))
	req.Header.Set("Content-Range", contentRange)

	// Execute request with retry logic for 5xx errors
	maxRetries := 3
	var resp *http.Response

	client := &http.Client{
		Timeout: 5 * time.Minute,
	}

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			waitTime := time.Duration(attempt) * 2 * time.Second
			fmt.Printf("Retrying chunk %d/%d after %v (attempt %d/%d)\n",
				chunkNum, totalChunks, waitTime, attempt+1, maxRetries)
			time.Sleep(waitTime)
		}

		resp, err = client.Do(req)
		if err != nil {
			if attempt == maxRetries-1 {
				return fmt.Errorf("failed to execute request after %d attempts: %w", maxRetries, err)
			}
			continue
		}

		// If we get a 5xx error, retry
		if resp.StatusCode >= 500 {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			if attempt == maxRetries-1 {
				return fmt.Errorf("server error (%d) after %d attempts: %s", resp.StatusCode, maxRetries, string(body))
			}
			continue
		}

		// Success or non-retryable error, break out of retry loop
		break
	}

	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle response based on status code
	switch resp.StatusCode {
	case http.StatusCreated: // 201 - Final chunk uploaded
		if chunkNum != totalChunks {
			return fmt.Errorf("received 201 status but not on final chunk (chunk %d/%d)", chunkNum, totalChunks)
		}
		return nil
	case http.StatusPartialContent: // 206 - Chunk uploaded successfully, more to go
		if chunkNum == totalChunks {
			return fmt.Errorf("received 206 status on final chunk (chunk %d/%d)", chunkNum, totalChunks)
		}
		return nil
	case http.StatusBadRequest: // 400
		return fmt.Errorf("bad request (400): %s", string(body))
	case http.StatusForbidden: // 403
		return fmt.Errorf("upload URL expired (403): %s", string(body))
	case http.StatusNotFound: // 404
		return fmt.Errorf("upload task not found (404): %s", string(body))
	case http.StatusRequestedRangeNotSatisfiable: // 416
		return fmt.Errorf("content range mismatch (416): %s", string(body))
	default:
		return fmt.Errorf("unexpected status code (%d): %s", resp.StatusCode, string(body))
	}
}

/*
func UploadFile(
	uploadUrl string,
	contentRange string,
	contentLength int,
	contentType string,
	filePath string,
) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Parse Content-Range to get the byte range to read
	// Format: "bytes {FIRST_BYTE}-{LAST_BYTE}/{TOTAL_BYTE_LENGTH}"
	var firstByte, lastByte, totalBytes int64
	_, err = fmt.Sscanf(contentRange, "bytes %d-%d/%d", &firstByte, &lastByte, &totalBytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse content range: %w", err)
	}

	// Calculate the actual chunk size to read
	chunkSize := lastByte - firstByte + 1

	// Validate that contentLength matches the calculated chunk size
	if int64(contentLength) != chunkSize {
		return "", fmt.Errorf("contentLength (%d) does not match calculated chunk size (%d)", contentLength, chunkSize)
	}

	// Seek to the starting position
	_, err = file.Seek(firstByte, 0)
	if err != nil {
		return "", fmt.Errorf("failed to seek file: %w", err)
	}

	// Read the chunk data
	chunkData := make([]byte, chunkSize)
	bytesRead, err := io.ReadFull(file, chunkData)
	if err != nil && err != io.ErrUnexpectedEOF {
		return "", fmt.Errorf("failed to read chunk: %w", err)
	}

	// Verify we read the expected amount
	if int64(bytesRead) != chunkSize {
		return "", fmt.Errorf("read %d bytes but expected %d", bytesRead, chunkSize)
	}

	// Create the HTTP request
	req, err := http.NewRequest("PUT", uploadUrl, bytes.NewReader(chunkData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Length", strconv.Itoa(contentLength))
	req.Header.Set("Content-Range", contentRange)

	// Execute the request
	client := &http.Client{
		Timeout: 5 * time.Minute, // Adjust timeout as needed
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle response based on status code
	switch resp.StatusCode {
	case http.StatusCreated: // 201 - All parts uploaded
		return string(body), nil
	case http.StatusPartialContent: // 206 - Chunk uploaded successfully, more to go
		return string(body), nil
	case http.StatusBadRequest: // 400
		return "", fmt.Errorf("bad request (400): %s", string(body))
	case http.StatusForbidden: // 403
		return "", fmt.Errorf("upload URL expired (403): %s", string(body))
	case http.StatusNotFound: // 404
		return "", fmt.Errorf("upload task not found (404): %s", string(body))
	case http.StatusRequestedRangeNotSatisfiable: // 416
		return "", fmt.Errorf("content range mismatch (416): %s", string(body))
	default:
		if resp.StatusCode >= 500 {
			// 5xx errors should be retried according to API docs
			return "", fmt.Errorf("server error (%d), should retry: %s", resp.StatusCode, string(body))
		}
		return "", fmt.Errorf("unexpected status code (%d): %s", resp.StatusCode, string(body))
	}
}

*/

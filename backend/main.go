package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type TokenRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
}

type FileRequest struct {
	ID              uuid.UUID `json:"id"`
	FilePathVideo   string    `json:"filepath_video"`
	FilePathContext string    `json:"filepath_context"`
	Status          string    `json:"status"`
}

type UploadRequest struct {
	Title          string `json:"title"`
	PrivacyLevel   string `json:"privacy_level"`
	FilePath       string `json:"file_path"`
	FileSize       int64  `json:"file_size"`
	CoverTimestamp int    `json:"cover_timestamp"`
}

func waitForAPI(apiURL string, maxRetries int) error {
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get(apiURL + "/health")
		if err == nil && resp.StatusCode == 200 {
			log.Println("API is ready!")
			return nil
		}
		log.Printf("Waiting for API... (attempt %d/%d URL %s)", i+1, maxRetries, apiURL)
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("API not ready after %d attempts", maxRetries)
}

func main() {

	var apiURL string
	var basePath string

	/*
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}*/

	mode := os.Getenv("MODE")

	fmt.Println(mode)

	if mode == "DEV" {
		apiURL = "http://localhost:8080"
		basePath = "/home/marcel/dev/scripts/go/backend"
	} else if mode == "PROD" {
		apiURL = "http://api:8080"
		basePath = "/home/marcel/app/backend"
	}

	// Warte bis API bereit ist
	if err := waitForAPI(apiURL, 10); err != nil {
		log.Fatal(err)
	}

	// Bind Exec Params
	filePathVideo := flag.String("video", "", "Path to video file")
	filePathContext := flag.String("context", "", "Path to context file")

	flag.Parse()

	video := *filePathVideo

	if *filePathVideo == "" {
		fmt.Println("No filepath for video provided")
		return
	}

	if *filePathContext == "" {
		fmt.Println("No filepath for context provided")
		return
	}

	title, err := ReadTitleFromContext(*filePathContext)
	if err != nil {
		fmt.Println("Error Extracting Name from JSON:", err)
		return
	}

	uuid := uuid.New()
	uuidString := uuid.String()

	fileName := basePath + "/videos/uploading/" + uuidString + "_UPLOADING.mp4"

	// Rename File
	err = os.Rename(video, fileName)
	if err != nil {
		fmt.Println("Error renaming file:", err)
		return
	}

	payloadFile := FileRequest{
		ID:              uuid,
		FilePathVideo:   fileName,
		FilePathContext: *filePathContext,
	}

	SentFile(payloadFile, apiURL)

	originalRefreshToken := "rft.7yekSfYUqyhHt7f6Inz3wkJ9ErZZ0lZkbuFrejf5n0KuKYXZcL13x3GqTuZV!4736.e1"

	// Renew Access Token and get Expiry in Seconds
	accessToken, expiresIn, err := RenewAccessToken(originalRefreshToken)
	if err != nil {
		panic(err)
	}

	// Load Time Zone
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		panic(err)
	}

	// Calculate expires_at using TikTok's expires_in
	expiresAt := time.Now().In(loc).Add(time.Duration(expiresIn) * time.Second).Format(time.RFC3339)

	payloadToken := TokenRequest{
		AccessToken:  accessToken,
		RefreshToken: originalRefreshToken, // optionally replace with new refresh token if returned
		ExpiresAt:    expiresAt,
	}

	SentToken(payloadToken, apiURL)

	fileSize, _, err := GetFileSize(fileName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	payloadUpload := UploadRequest{
		Title:          title,
		PrivacyLevel:   "SELF_ONLY",
		FilePath:       fileName,
		FileSize:       fileSize,
		CoverTimestamp: 1000,
	}

	SentUpload(payloadUpload, apiURL)

	// Create Upload URL for File

	uploadUrl, err := CreateUploadURL(

		title,
		"SELF_ONLY",
		fileName,
		1000,
		accessToken,
		originalRefreshToken,
	)

	if err != nil {
		fmt.Println("Error:", err)

		filePathFailed := basePath + "/videos/failed/" + uuidString + "_FAILED.mp4"
		payloadFile = FileRequest{
			ID:              uuid,
			FilePathVideo:   filePathFailed,
			FilePathContext: *filePathContext,
			Status:          "FAILED",
		}

		UpdateFile(payloadFile, apiURL)

		err = os.Rename(fileName, filePathFailed)
		if err != nil {
			fmt.Println("Error renaming file:", err)
			return
		}

	}

	contentType := "video/mp4"
	err = UploadFileComplete(uploadUrl, fileName, fileSize, contentType)
	if err != nil {

		filePathFailed2 := basePath + "/videos/failed/" + uuidString + "_FAILED.mp4"
		payloadFile = FileRequest{
			ID:              uuid,
			FilePathVideo:   filePathFailed2,
			FilePathContext: *filePathContext,
			Status:          "FAILED",
		}

		UpdateFile(payloadFile, apiURL)

		err = os.Rename(fileName, filePathFailed2)
		if err != nil {
			fmt.Println("Error renaming file:", err)
			return
		}

		fmt.Printf("Error: %v\n", err)
	}

	filePathUploaded := basePath + "/videos/uploaded/" + uuidString + "_UPLOADED.mp4"

	payloadFile = FileRequest{
		ID:              uuid,
		FilePathVideo:   filePathUploaded,
		FilePathContext: *filePathContext,
		Status:          "UPLOADED",
	}

	UpdateFile(payloadFile, apiURL)

	// Rename File
	err = os.Rename(fileName, filePathUploaded)

	if err != nil {
		fmt.Println("Error renaming file:", err)
		return
	}

}

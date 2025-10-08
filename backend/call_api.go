package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func SentUpload(payload UploadRequest, apiURL string) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(
		apiURL+"/api/upload/create",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	return nil
}

func SentToken(payload TokenRequest, apiURL string) *http.Response {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(
		apiURL+"/api/token/create",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	return resp
}

func SentFile(payload FileRequest, apiURL string) *http.Response {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(
		apiURL+"/api/file/create",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	return resp
}

func UpdateFile(payload FileRequest, apiURL string) *http.Response {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(
		http.MethodPut,
		apiURL+"/api/file/update",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		panic(err)
	}

	// Set content type header
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	return resp
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func SentUpload(payload UploadRequest) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(
		"http://127.0.0.1:8080/api/upload/create",
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

func SentToken(payload TokenRequest) *http.Response {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

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

	return resp
}

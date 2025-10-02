package main

import (
	"log"

	"github.com/grundigdev/club/models"
	"github.com/grundigdev/club/shared"
)

func main() {
	db, err := shared.NewPostgres()
	if err != nil {
		panic(err)
	}

	// Support for UUID
	err = db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error
	if err != nil {
		panic(err)
	}

	// AutoMigrate
	err = db.AutoMigrate(&models.TokenModel{}, &models.UploadModel{})
	if err != nil {
		panic(err)
	}

	log.Println("MIGRATION SUCCESSFUL")
}

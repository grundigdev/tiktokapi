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

	// Check if enum type 'trade_status' exists, if not, create it

	err = db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error
	if err != nil {
		panic(err)
	}

	// AutoMigrate for the other tables
	err = db.AutoMigrate(&models.TokenModel{})
	if err != nil {
		panic(err)
	}

	log.Println("MIGRATION SUCCESSFUL")
}

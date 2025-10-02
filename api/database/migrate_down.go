package main

import (
	"log"

	"github.com/grundigdev/club/shared"
)

func main() {
	db, err := shared.NewPostgres()
	if err != nil {
		panic(err)
	}

	// Drop all tables in the database
	allTables := []string{}
	err = db.Raw(`
		SELECT tablename 
		FROM pg_tables 
		WHERE schemaname = 'public';
	`).Scan(&allTables).Error
	if err != nil {
		panic(err)
	}

	for _, table := range allTables {
		log.Printf("Dropping table: %s", table)
		if err := db.Migrator().DropTable(table); err != nil {
			panic(err)
		}
	}

	log.Println("ALL TABLES HAS BEEN DROPPED SUCCESSFULLY")
}

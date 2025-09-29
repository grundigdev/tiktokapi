package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/grundigdev/club/requests"
	"github.com/grundigdev/club/services"
	"github.com/grundigdev/club/shared"
)

func main() {
	// Initialize database connection
	db, err := shared.NewPostgres()
	if err != nil {
		panic(err)
	}
	absenceService := services.AbsenceService{
		DB: db,
	}

	// Generate 1000 random absences
	for i := 0; i < 1000; i++ {
		// Define the start and end of the year 2024
		start := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)

		// Calculate the total seconds in the year 2024
		seconds := int(end.Sub(start).Seconds())

		// Generate a random start date
		randomSecondsBig, err := rand.Int(rand.Reader, big.NewInt(int64(seconds)))
		if err != nil {
			panic("Failed to generate random seconds: " + err.Error())
		}
		randomSeconds := int(randomSecondsBig.Int64())
		randomStartDate := start.Add(time.Duration(randomSeconds) * time.Second)

		// Generate a random duration for the absence (1 to 7 days)
		randomDaysBig, err := rand.Int(rand.Reader, big.NewInt(7))
		if err != nil {
			panic("Failed to generate random duration: " + err.Error())
		}
		randomDurationDays := int(randomDaysBig.Int64()) + 1
		randomEndDate := randomStartDate.Add(time.Duration(randomDurationDays) * 24 * time.Hour)

		// Generate random Absence Type ID (1 to 5)
		randomAbsenceTypeIDBig, err := rand.Int(rand.Reader, big.NewInt(5))
		if err != nil {
			panic("Failed to generate random Absence Type ID: " + err.Error())
		}
		randomAbsenceTypeID := uint(randomAbsenceTypeIDBig.Int64()) + 1

		// Generate random UserID (1 to 3)
		randomUserIDBig, err := rand.Int(rand.Reader, big.NewInt(3))
		if err != nil {
			panic("Failed to generate random UserID: " + err.Error())
		}
		randomUserID := uint(randomUserIDBig.Int64()) + 1

		// Create the absence request
		absenceRequest := &requests.GenerateAbsenceRequest{
			TypeID:    randomAbsenceTypeID,
			StartDate: randomStartDate,
			EndDate:   randomEndDate,
			UserID:    randomUserID,
		}

		// Call the service to create the absence
		_, err = absenceService.GenerateAbsence(absenceRequest)
		if err != nil {
			fmt.Printf("Failed to create absence: %v\n", err)
			continue
		}

		fmt.Printf("Absence Type %d for User %d from %s to %s created\n",
			absenceRequest.TypeID,
			absenceRequest.UserID,
			absenceRequest.StartDate.Format("2006-01-02"),
			absenceRequest.EndDate.Format("2006-01-02"))
	}
}

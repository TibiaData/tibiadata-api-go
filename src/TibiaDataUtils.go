package main

import (
	"log"
	"time"

	"golang.org/x/text/unicode/norm"
)

// TibiadataDatetimeV3 func
func TibiadataDatetimeV3(date string) string {
	//TODO: Normalization needs to happen above this layer
	date = norm.NFKC.String(date)

	var returnDate time.Time
	var err error

	// If statement to determine if date string is filled or empty
	if date == "" {
		// The string that should be returned is the current timestamp
		returnDate = time.Now()
	} else {
		// Parse: Jan 02 2007, 19:20:30 CET
		returnDate, err = time.Parse("Jan 02 2006, 15:04:05 MST", date)

		if err != nil {
			log.Println(err)
		}
	}

	// Return of formatted date and time string to functions..
	return returnDate.UTC().Format(time.RFC3339)
}

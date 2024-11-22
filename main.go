package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading.env file: %v\n", err)
	}

	log.Printf("start scraping")
	fecultyList := fetchAllFaculties()

	// writeFacultyToCSV("faculties.csv", fecultyList)
	generateReadme(fecultyList)

}

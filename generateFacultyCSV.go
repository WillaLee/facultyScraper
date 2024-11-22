package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

// writeFacultyToCSV writes a slice of Faculty objects to a CSV file
func writeFacultyToCSV(filename string, facultyList []Faculty) error {
	// Open the file for writing
	log.Print("start writing faculties to CSV")
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write BOM for UTF-8 encoding
	file.Write([]byte{0xEF, 0xBB, 0xBF})

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header
	header := []string{"Name", "Link", "Position Titles", "Locations", "Research Interests"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write the data
	for _, faculty := range facultyList {
		record := []string{
			faculty.Name,
			faculty.Link,
			faculty.PosTitles,
			fmt.Sprintf("%v", faculty.Locations), // Convert slice to string
			fmt.Sprintf("%v", faculty.ResearchInterests), // Convert slice to string
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

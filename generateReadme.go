package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Function to generate a README.md file with a faculty table
func generateReadme(faculties []Faculty) error {
	// Open file for writing
	file, err := os.Create("README.md")
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	description := `This project scrapes faculty information from
	 the [Khoury College of Computer Sciences](https://www.khoury.northeastern.edu/about/people/) at Northeastern University.
	  It collects details including the faculty's name, title, location, research interests, and a link to their profile. 
	  You can also download the [CSV file](https://github.com/WillaLee/facultyScraper/blob/main/faculties.csv) for easy access and analysis.` + "\n\n" +
		`> **Note:** There may be some errors in the table because the faculty pages may have slightly different formats. For the most accurate information, please visit the faculty pages directly.` + "\n\n"

	_, err = file.WriteString(description)
	if err != nil {
		return fmt.Errorf("error writing description: %v", err)
	}

	// Write header
	_, err = file.WriteString("# Faculty Directory\n\n")
	if err != nil {
		return fmt.Errorf("error writing header: %v", err)
	}

	// Write table header
	tableHeader := "| Name | Position | Location | Research Interests |\n" +
		"|------|-----------|-----------|-------------------|\n"
	_, err = file.WriteString(tableHeader)
	if err != nil {
		return fmt.Errorf("error writing table header: %v", err)
	}

	// Write table rows
	for _, faculty := range faculties {
		// Create markdown link for name
		nameLink := fmt.Sprintf("[%s](%s)", faculty.Name, faculty.Link)

		// Format locations as comma-separated string
		locations := strings.Join(faculty.Locations, ", ")

		// Format research interests as comma-separated string
		interests := strings.Join(faculty.ResearchInterests, ", ")

		// Replace any pipe characters in text fields to avoid breaking the table
		nameLink = strings.ReplaceAll(nameLink, "|", "\\|")
		positions := strings.ReplaceAll(faculty.PosTitles, "|", "\\|")
		locations = strings.ReplaceAll(locations, "|", "\\|")
		interests = strings.ReplaceAll(interests, "|", "\\|")

		// Create table row
		row := fmt.Sprintf("| %s | %s | %s | %s |\n",
			nameLink,
			positions,
			locations,
			interests)

		_, err = file.WriteString(row)
		if err != nil {
			return fmt.Errorf("error writing row: %v", err)
		}
	}

	// Add footer with generation timestamp
	footer := fmt.Sprintf("\n\n*Last updated: %s*\n", time.Now().Format("January 2, 2006"))
	_, err = file.WriteString(footer)
	if err != nil {
		return fmt.Errorf("error writing footer: %v", err)
	}

	return nil
}

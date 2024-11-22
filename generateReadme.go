package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

// Function to generate a README.md file with a faculty table
func generateReadme(faculties []Faculty) error {
	// Define the table structure for the markdown format
	const readmeTemplate = `
# Faculty Directory

This is the list of faculties with their details.

| Name           | Position Title    | Locations | Research Interests | Link                               |
|----------------|-------------------|-----------|---------------------|------------------------------------|
{{range .}}
| {{.Name}}      | {{.PosTitles}}     | {{join .Locations ", "}} | {{join .ResearchInterests ", "}} | [Link]({{.Link}}) |
{{end}}

`

	// Create a new file or open for writing
	file, err := os.Create("README.md")
	if err != nil {
		return fmt.Errorf("unable to create README.md: %v", err)
	}
	defer file.Close()

	// Prepare a template for the table
	tmpl, err := template.New("readme").Funcs(template.FuncMap{
		"join": func(slice []string, sep string) string {
			return strings.Join(slice, sep)
		},
	}).Parse(readmeTemplate)
	if err != nil {
		return fmt.Errorf("unable to parse template: %v", err)
	}

	// Write the content to the file using the template
	err = tmpl.Execute(file, faculties)
	if err != nil {
		return fmt.Errorf("unable to execute template: %v", err)
	}

	log.Println("README.md generated successfully.")
	return nil
}

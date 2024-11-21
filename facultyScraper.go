package main

import (
	"log"
	"encoding/json"

	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Faculty struct {
	Title      string   `json:"title"`
	Link       string   `json:"link"`
	PosTitles  string   `json:"pos_titles"`
	Locations  []string `json:"locations"`
}

func fetchAllFacultyPageAJAX(page int) ([]Faculty, error){
	// Define the URL for the AJAX request
	url := "https://www.khoury.northeastern.edu/wp-admin/admin-ajax.php"

	// Prepare the POST request payload
	startNumber := page * 20
	payload := fmt.Sprintf("ptype=people&research_area=0&location=0&role=0&maincat=faculty&offset=%d&action=people_filter_ajax", startNumber)

	// Define the necessary headers for the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
		return nil, err;
	}

	// Set the headers
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Mobile Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Referer", "https://www.khoury.northeastern.edu/about/people/")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
		return nil, err;
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
		return nil, err;
	}

	// parse response body into a slice of faculties
	var faculties []Faculty
	err = json.Unmarshal(body, &faculties)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
		return nil, err;
	}

	return faculties, nil
}

func fetchAllFaculties() []Faculty{
	var allFaculties []Faculty
	for i := 0; i <= 29; i++ {
		facultiesInPage, err := fetchAllFacultyPageAJAX(i)
		if err != nil {
			log.Fatalf("Error fetching page %d: %v", i, err)
		} else {
			allFaculties = append(allFaculties, facultiesInPage...)
		}
	}

	return allFaculties;
}
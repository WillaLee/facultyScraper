package main

import (
	"encoding/json"
	"log"
	"sort"
	"sync"

	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Faculty struct {
	Name              string   `json:"title"`
	Link              string   `json:"link"`
	PosTitles         string   `json:"pos_titles"`
	Locations         []string `json:"locations"`
	ResearchInterests []string
}

// ByName implements sort.Interface for sorting Faculty slices by Name.
type ByName []Faculty

// Len returns the number of elements in the collection.
func (a ByName) Len() int {
	return len(a)
}

// Less reports whether the element with index i should sort before the element with index j.
func (a ByName) Less(i, j int) bool {
	return a[i].Name < a[j].Name
}

// Swap swaps the elements with indexes i and j.
func (a ByName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func fetchAllFacultyPageAJAX(page int) ([]Faculty, error) {
	// Define the URL for the AJAX request
	url := "https://www.khoury.northeastern.edu/wp-admin/admin-ajax.php"

	// Prepare the POST request payload
	startNumber := page * 20
	payload := fmt.Sprintf("ptype=people&research_area=0&location=0&role=0&maincat=faculty&offset=%d&action=people_filter_ajax", startNumber)

	// Define the necessary headers for the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
		return nil, err
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
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
		return nil, err
	}

	// parse response body into a slice of faculties
	var faculties []Faculty
	err = json.Unmarshal(body, &faculties)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
		return nil, err
	}

	return faculties, nil
}

// fetch all faculties
func fetchAllFaculties() []Faculty {
	var (
		allFaculties  []Faculty
		mu            sync.Mutex
		wg            sync.WaitGroup
		pageLimit     = 8 // Limit for concurrent goroutines per page
		interestLimit = 5 // Limit for concurrent goroutines per faculty
	)

	// Channel to limit concurrency for fetching faculty pages
	pageSemaphore := make(chan struct{}, pageLimit)

	// Channel to limit concurrency for fetching research interests
	interestSemaphore := make(chan struct{}, interestLimit)

	for i := 0; i <= 29; i++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			pageSemaphore <- struct{}{}
			defer func() { <-pageSemaphore }()

			facultiesInPage, err := fetchAllFacultyPageAJAX(page)
			if err != nil {
				log.Printf("Error fetching page %d: %v\n", page, err)
				return
			}

			var pageWg sync.WaitGroup
			for _, faculty := range facultiesInPage {
				pageWg.Add(1)
				go func(f Faculty) {
					defer pageWg.Done()
					interestSemaphore <- struct{}{}
					defer func() { <-interestSemaphore }()

					researchInterests, err := fetchResearchInterests(f.Link)
					if err != nil {
						log.Printf("Error fetching research interests for %s: %v\n", f.Name, err)
						return
					}
					f.ResearchInterests = researchInterests

					if len(researchInterests) > 0 {
						// printFacultyDetails(f)
						mu.Lock()
						allFaculties = append(allFaculties, f)
						mu.Unlock()
					}
				}(faculty)
			}
			pageWg.Wait()
		}(i)
	}

	wg.Wait()

	log.Print("start sort")
	sort.Sort(ByName(allFaculties))
	return allFaculties
}

func printFacultyDetails(f Faculty) {
	fmt.Println("Faculty Details:")
	fmt.Println("-----------------")
	fmt.Printf("Name: %s\n", f.Name)
	fmt.Printf("Link: %s\n", f.Link)
	fmt.Printf("Position Titles: %s\n", f.PosTitles)

	fmt.Print("Locations: ")
	if len(f.Locations) > 0 {
		fmt.Println()
		for _, location := range f.Locations {
			fmt.Printf("  - %s\n", location)
		}
	} else {
		fmt.Println("None")
	}

	fmt.Print("Research Interests: ")
	if len(f.ResearchInterests) > 0 {
		fmt.Println()
		for _, interest := range f.ResearchInterests {
			fmt.Printf("  - %s\n", interest)
		}
	} else {
		fmt.Println("None")
	}
	fmt.Println()
}

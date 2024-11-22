package main

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// fetchResearchInterests fetches the HTML file and extracts the "Research Interests" list
func fetchResearchInterests(url string) ([]string, error) {
	// Fetch the HTML file
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	// Check for HTTP response errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Select the "Research Interests" section and extract the list items
	var interests []string
	doc.Find(`div.copy.paragraph[data-swiftype-index="true"][data-swiftype-name="resultDescription"][data-swiftype-type="text"]`).Each(func(i int, s *goquery.Selection) {
		s.Find("h2").Each(func(i int, h *goquery.Selection) {
			if h.Text() == "Research Interests" || h.Text() == "Research Areas" {
				h.Next().Find("li").Each(func(j int, li *goquery.Selection) {
					interests = append(interests, li.Text())
				})
			}
		})
	})

	return interests, nil
}

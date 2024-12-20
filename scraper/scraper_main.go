package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"

	"quotes/config"
	"quotes/connection"
)

const (
	maxRetries      = 3
	retryDelay      = 2 * time.Second
	rateLimitPeriod = 500 * time.Millisecond
)

// ScrapeQuotes fetches quotes from the specified BASE_URL up to the desired count
func ScrapeQuotes(baseUrl string, count int) ([]connection.Quote, error) {
	var quotes []connection.Quote
	page := 1

	rateLimiter := time.Tick(rateLimitPeriod)

	for len(quotes) < count {
		url := fmt.Sprintf("%s/page/%d/", baseUrl, page)
		log.Printf("Scraping URL: %s", url)

		<-rateLimiter

		resp, err := fetchWithRetry(url)
		if err != nil {
			log.Printf("Failed to fetch page %d: %v. Skipping to next page.", page, err)
			page++
			continue
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("Received non-OK HTTP status %d for page %d", resp.StatusCode, page)
			resp.Body.Close()
			break
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Printf("Failed to parse HTML for page %d: %v", page, err)
			page++
			continue
		}

		doc.Find("div.quote").Each(func(i int, s *goquery.Selection) {
			text := s.Find("span.text").Text()
			author := s.Find("small.author").Text()
			text = cleanQuoteText(text)

			quotes = append(quotes, connection.Quote{
				Text:   text,
				Author: author,
			})
		})

		log.Printf("Scraped %d quotes from page %d", len(quotes), page)
		page++
	}

	return quotes, nil
}

// fetchWithRetry performs an HTTP GET request with retry logic
func fetchWithRetry(url string) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i < maxRetries; i++ {
		resp, err = http.Get(url)
		if err == nil {
			return resp, nil
		}

		log.Printf("Retrying... attempt %d/%d for URL: %s", i+1, maxRetries, url)
		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("failed to fetch URL %s after %d attempts: %v", url, maxRetries, err)
}

// cleanQuoteText removes the surrounding quotes from the quote text
func cleanQuoteText(text string) string {
	runes := []rune(text)
	if len(runes) >= 2 && runes[0] == '“' && runes[len(runes)-1] == '”' {
		return string(runes[1 : len(runes)-1])
	}
	return text
}

func main() {
	connection.InitCache()
	desiredCount := config.DefaultCount
	log.Printf("Starting scraper to fetch %d quotes", desiredCount)

	quotes, err := ScrapeQuotes(config.BaseUrl, desiredCount)
	if err != nil {
		log.Fatalf("Error scraping quotes: %v", err)
	}

	log.Printf("Successfully scraped %d quotes", len(quotes))

	err = connection.StoreQuotes(quotes)
	if err != nil {
		log.Fatalf("Error storing quotes in Redis: %v", err)
	}

	log.Println("Quotes successfully stored in Redis")
}

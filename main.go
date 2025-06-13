package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

func ExtractURLs(fileContent string) ([]string, error) {
	urlPattern := `(https?://[^\s\)"'<>]+)|\[.*?\]\((.*?)\)`
	re, err := regexp.Compile(urlPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %w", err)
	}
	var urls []string

	matches := re.FindAllStringSubmatch(fileContent, -1)

	for _, match := range matches {
		if match[1] != "" && hasHTTPPrefix(match[1]) {
			urls = append(urls, match[1])
		}
		if match[2] != "" && hasHTTPPrefix(match[2]) {
			urls = append(urls, match[2])
		}
	}
	seen := make(map[string]bool)
	uniqueURLs := []string{}
	for _, url := range urls {
		if !seen[url] {
			seen[url] = true
			uniqueURLs = append(uniqueURLs, url)
		}
	}
	return uniqueURLs, nil
}

func hasHTTPPrefix(s string) bool {
	return len(s) > 7 && (s[:7] == "http://" || (len(s) > 8 && s[:8] == "https://"))
}

func CheckURL(client *http.Client, jobs <-chan string, wg *sync.WaitGroup, successResults chan<- string, failureResults chan<- string) {
	defer wg.Done()
	for url := range jobs {
		_, err := client.Head(url)
		if err != nil {
			// If HEAD request fails, try GET request
			_, err := client.Get(url)
			if err != nil {
				failureResults <- url
				continue
			}
		}
		successResults <- url
	}
}

func main() {
	filePath := os.Getenv("INPUT_FILE_PATH")
	workerStr := os.Getenv("INPUT_CONCURRENT_WORKERS")
	timeoutSeconds := os.Getenv("INPUT_TIMEOUT_SECONDS")
	maxConcurrentWorkers := 30
	defaultTimeoutSeconds := 5

	if filePath == "" {
		fmt.Println("Error: file_path value is required")
		os.Exit(1)
	}

	if workerStr != "" {
		var err error
		maxConcurrentWorkers, err = strconv.Atoi(workerStr)
		if err != nil || maxConcurrentWorkers < 1 {
			fmt.Printf("⚠️ Invalid concurrent_workers: %s\n", workerStr)
			fmt.Println("ℹ️ Using default value of 30 for concurrent workers.")
			maxConcurrentWorkers = 30
		}
	}
	if timeoutSeconds != "" {
		var err error
		defaultTimeoutSeconds, err = strconv.Atoi(timeoutSeconds)
		if err != nil || defaultTimeoutSeconds < 1 {
			fmt.Printf("⚠️ Invalid timeout_seconds: %s\n", timeoutSeconds)
			fmt.Println("ℹ️ Using default value of 5 seconds for timeout.")
			defaultTimeoutSeconds = 5
		}
	}

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}
	fileContent := string(fileBytes)
	urls, err := ExtractURLs(fileContent)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	var wg sync.WaitGroup
	successResults := make(chan string, len(urls))
	failureResults := make(chan string, len(urls))
	jobs := make(chan string, len(urls))

	client := &http.Client{
		Timeout: time.Duration(defaultTimeoutSeconds) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	for worker := 1; worker <= maxConcurrentWorkers; worker++ {
		wg.Add(1)
		go CheckURL(client, jobs, &wg, successResults, failureResults)
	}

	for _, url := range urls {
		jobs <- url
	}
	close(jobs) // close the channel when there's no more to send.

	go func() {
		wg.Wait()             // Wait for all Go routines to complete
		close(successResults) // Close the results channel after all Go routines finish
		close(failureResults) // Close the results channel after all Go routines finish
	}()

	fmt.Println("✅ Working URLs:")
	for result := range successResults {
		fmt.Printf("- %s\n", result)
	}

	fmt.Println("\n\n❌ Invalid URLs:")
	var invalidURLs []string
	for result := range failureResults {
		invalidURLs = append(invalidURLs, result)
		fmt.Printf("- %s\n", result)
	}

	if len(invalidURLs) != 0 {
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	serverURL := os.Getenv("HELLO_SERVER_URL")
	if serverURL == "" {
		log.Fatal("HELLO_SERVER_URL env variable not set")
	}

	// Set up dummy HTTP server to satisfy Cloud Run requirements
	go func() {
		http.HandleFunc("/dummy", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Main logic
	numThreads := 100
	var wg sync.WaitGroup
	wg.Add(numThreads)

	startTime := time.Now()
	// Start 100 threads
	for i := 0; i < numThreads; i++ {
		go func() {
			defer wg.Done()
			makeGetRequest(serverURL)
		}()
	}
	wg.Wait()
	duration := time.Since(startTime)
	fmt.Printf("Took %v to make %d requests\n", duration, numThreads)
}

func makeGetRequest(url string) {
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("%s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s", bodyBytes)
}

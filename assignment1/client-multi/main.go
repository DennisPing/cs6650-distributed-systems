package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/DennisPing/cs6650-distributed-systems/assignment1/client-single/log"
	api "github.com/DennisPing/twinder-sdk-go"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// 	maxWorkers  = 100
	numRequests = 500_000
	maxWorkers  = 10
)

func main() {
	serverURL := os.Getenv("SERVER_URL")
	if serverURL == "" {
		log.Logger.Fatal().Msg("SERVER_URL env variable not set")
	}

	// Set up dummy HTTP server to satisfy Cloud Run requirements
	go func() {
		http.HandleFunc("/dummy", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})
		log.Logger.Fatal().Msg(http.ListenAndServe(":8081", nil).Error())
	}()

	ctx := context.Background()

	// Initialize RNG client pool
	RngClientPool := make(chan *RngClient, maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		RngClientPool <- NewRngClient(serverURL)
	}

	var wg sync.WaitGroup
	for i := 0; i < numRequests; i++ {
		direction := "left"
		if rand.Intn(2) == 1 { // Flip a coin
			direction = "right"
		}
		wg.Add(1)

		go func() {
			defer wg.Done()
			rngClient := <-RngClientPool // Get an RNG client from the pool
			defer func() {
				RngClientPool <- rngClient // Release the slot so that other goroutines can acquire slot
			}()
			swipeLeftOrRight(ctx, rngClient, direction)
		}()
	}
	wg.Wait()
	fmt.Println("Done!")
}

func swipeLeftOrRight(ctx context.Context, client *RngClient, direction string) {
	reqBody := api.SwipeDetails{
		Swiper:  strconv.Itoa(randInt(client.rng, 1, 5000)),
		Swipee:  strconv.Itoa(randInt(client.rng, 1, 1_000_000)),
		Comment: randComment(client.rng, 256),
	}

	// Send POST request
	resp, err := client.apiClient.SwipeApi.Swipe(ctx, reqBody, direction)
	if err != nil {
		log.Logger.Error().Interface("swaggerError", err)
		return
	}

	// StatusCode should be 200 or 201, else log warn
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		log.Logger.Debug().Msg(resp.Status)
	} else {
		log.Logger.Warn().Msg(resp.Status)
	}
}

// Each goroutine client should pass in their own RNG
func randInt(rng *rand.Rand, start, stop int) int {
	return rng.Intn(stop) + start
}

// Each goroutine client should pass in their own RNG
func randComment(rng *rand.Rand, length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rng.Intn(len(charset))]
	}
	return string(b)
}

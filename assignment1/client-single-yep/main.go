package main

import (
	"context"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/DennisPing/cs6650-distributed-systems/assignment1/client-single/log"
	api "github.com/DennisPing/twinder-client-sdk"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

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

	cfg := api.NewConfiguration()
	apiClient := api.NewAPIClient(cfg)
	apiClient.ChangeBasePath(serverURL)

	rand.Seed(time.Now().UnixNano())
	ctx := context.Background()

	swipeLeftOrRight(ctx, apiClient, "left")
	swipeLeftOrRight(ctx, apiClient, "right")
}

func swipeLeftOrRight(ctx context.Context, client *api.APIClient, direction string) {
	reqBody := api.SwipeDetails{
		Swiper:  strconv.Itoa(randInt(1, 5000)),
		Swipee:  strconv.Itoa(randInt(1, 1000000)),
		Comment: randComment(256),
	}

	// Send POST request
	resp, err := client.SwipeApi.Swipe(ctx, reqBody, direction)
	if err != nil {
		log.Logger.Error().Interface("swaggerError", err)
		return
	}

	// StatusCode should be 200 or 201, else log warn
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		log.Logger.Info().Msg(resp.Status)
	} else {
		log.Logger.Warn().Msg(resp.Status)
	}
}

func randInt(start, stop int) int {
	return rand.Intn(stop) + start
}

func randComment(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

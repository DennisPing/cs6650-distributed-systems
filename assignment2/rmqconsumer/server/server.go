package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DennisPing/cs6650-distributed-systems/assignment2/consumer/logger"
	"github.com/DennisPing/cs6650-distributed-systems/assignment2/consumer/models"
	"github.com/DennisPing/cs6650-distributed-systems/assignment2/consumer/store"
	"github.com/go-chi/chi"
)

func Start(kvStore *store.SimpleStore) *http.Server {
	chiRouter := chi.NewRouter()

	// GET /health
	chiRouter.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		writeStatusResponse(w, http.StatusOK)
	})

	// GET /swipes?userId=1234
	chiRouter.Get("/swipes", func(w http.ResponseWriter, r *http.Request) {
		userId := r.URL.Query().Get("userId")
		stats, found := kvStore.GetUserStats(userId)
		if !found {
			http.Error(w, "userId not found", http.StatusNotFound)
			return
		}
		twinderLikes := &models.TwinderUserStats{
			UserId:   userId,
			Likes:    stats.Likes,
			Dislikes: stats.Dislikes,
		}
		writeJsonResponse(w, http.StatusOK, twinderLikes)
	})

	// GET /matches?userId=1234
	chiRouter.Get("/matches", func(w http.ResponseWriter, r *http.Request) {
		userId := r.URL.Query().Get("userId")
		stats, found := kvStore.GetUserStats(userId)
		if !found {
			http.Error(w, "userId not found", http.StatusNotFound)
			return
		}
		twinderMatches := &models.TwinderMatches{
			UserId:  userId,
			Matches: stats.Matches,
		}
		writeJsonResponse(w, http.StatusOK, twinderMatches)
	})

	// GET /matches/all
	chiRouter.Get("/swipes/all", func(w http.ResponseWriter, r *http.Request) {
		allStats := kvStore.GetAllUserStats()
		twinderUsersStats := &models.AllTwinderUserStats{
			UsersStats: []models.TwinderUserStats{},
		}

		for userId, stats := range allStats {
			userStats := models.TwinderUserStats{
				UserId:   userId,
				Likes:    stats.Likes,
				Dislikes: stats.Dislikes,
			}
			twinderUsersStats.UsersStats = append(twinderUsersStats.UsersStats, userStats)
		}
		writeJsonResponse(w, http.StatusOK, twinderUsersStats)
	})

	server := &http.Server{
		Addr:    ":8080", // Set this to 8081 or 8082 in local testing
		Handler: chiRouter,
	}

	go func() {
		logger.Info().Msg("starting HTTP server...")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("failed to start HTTP server")
		}
	}()

	return server
}

// Write a simple HTTP status to the response writer
func writeStatusResponse(w http.ResponseWriter, statusCode int) {
	logger.Info().Int("code", statusCode)
	w.WriteHeader(statusCode)
}

// Marshal and write a JSON response to the response writer
func writeJsonResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	respBytes, err := json.Marshal(payload)
	if err != nil {
		logger.Error().Int("code", http.StatusInternalServerError).Msg("error marshaling JSON response")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	logger.Info().Interface("send", payload).Send()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(respBytes)))
	w.WriteHeader(statusCode)
	w.Write(respBytes)
}

// Write an HTTP error to the response writer
func writeErrorResponse(w http.ResponseWriter, method string, statusCode int, message string) {
	logger.Warn().Str("method", method).Int("code", statusCode).Msg(message)
	http.Error(w, message, statusCode)
}

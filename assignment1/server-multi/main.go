package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/DennisPing/cs6650-distributed-systems/assignment1/server/log"
	"github.com/DennisPing/cs6650-distributed-systems/assignment1/server/models"
	"github.com/go-chi/chi"
)

const maxWorkers = 20

func main() {
	// Parse CLI flags

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := NewServer()

	log.Logger.Info().Msgf("starting the server on port %s", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), server)
	if err != nil {
		log.Logger.Fatal().Msg(err.Error())
	}
}

type Server struct {
	*chi.Mux
	workerPool chan func()
}

func NewServer() *Server {
	workerPool := make(chan func(), maxWorkers) // Buffered channel to hold tasks
	for i := 0; i < maxWorkers; i++ {
		go worker(workerPool)
	}
	s := &Server{
		Mux:        chi.NewRouter(),
		workerPool: workerPool,
	}
	s.Get("/", s.homeHandler)
	s.Post("/swipe/{leftorright}/", s.swipeHandler)
	return s
}

// Each worker listens on the taskQueue for a new task.
// If a task is available it executes the function.
func worker(taskQueue <-chan func()) {
	for task := range taskQueue {
		task()
	}
}

// Hello world endpoint for debugging purposes
func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello world!"))
}

// Handle swipe left or right
func (s *Server) swipeHandler(w http.ResponseWriter, r *http.Request) {
	leftorright := chi.URLParam(r, "leftorright")

	var wg sync.WaitGroup
	wg.Add(1)

	// Put this task into the worker pool channel if there's room, else block
	s.workerPool <- func() {
		defer wg.Done()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeErrorResponse(w, r.Method, http.StatusBadRequest, "bad request")
			return
		}

		var reqBody models.SwipeRequest
		err = json.Unmarshal(body, &reqBody)
		if err != nil {
			writeErrorResponse(w, r.Method, http.StatusBadRequest, "bad request")
			return
		}

		resp := models.SwipeResponse{
			Message: fmt.Sprintf("you swiped %s", leftorright),
		}
		// left and right do the same thing for now
		switch leftorright {
		case "left":
			writeJsonResponse(w, http.StatusCreated, resp)
		case "right":
			writeJsonResponse(w, http.StatusCreated, resp)
		default:
			writeErrorResponse(w, r.Method, http.StatusBadRequest, leftorright)
		}
	}
	wg.Wait()
}

// Marshal and write a JSON response to the response writer
func writeJsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	log.Logger.Info().Interface("send", data).Send()
	respJson, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "error marshaling JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(respJson)
}

// Write an HTTP error to the response writer
func writeErrorResponse(w http.ResponseWriter, method string, statusCode int, message string) {
	log.Logger.Error().
		Str("method", method).
		Int("code", statusCode).
		Msg(message)
	http.Error(w, message, statusCode)
}

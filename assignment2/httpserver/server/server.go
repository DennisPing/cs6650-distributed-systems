package server

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/DennisPing/cs6650-distributed-systems/assignment2/httpserver/log"
	"github.com/DennisPing/cs6650-distributed-systems/assignment2/httpserver/metrics"
	"github.com/DennisPing/cs6650-distributed-systems/assignment2/httpserver/models"
	"github.com/go-chi/chi"
	"github.com/wagslane/go-rabbitmq"
)

type Server struct {
	*chi.Mux
	*http.Server
	*metrics.Metrics
	Pub *rabbitmq.Publisher
}

func NewServer(address string, metrics *metrics.Metrics, publisher *rabbitmq.Publisher) *Server {
	chiRouter := chi.NewRouter()
	s := &Server{
		Mux: chiRouter,
		Server: &http.Server{
			Addr:    address,
			Handler: chiRouter,
		},
		Metrics: metrics,
		Pub:     publisher,
	}

	s.Get("/health", s.homeHandler)
	s.Post("/swipe/{leftorright}/", s.swipeHandler)
	return s
}

func (s *Server) Start() error {
	ticker := time.NewTicker(5 * time.Second)
	go func() { // Metrics goroutine
		for range ticker.C {
			err := s.SendMetrics()
			if err != nil {
				log.Error().Msgf("unable to send metrics to Axiom: %v", err)
			}
		}
	}()
	return s.ListenAndServe()
}

// Health endpoint
func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello world!"))
}

// Handle swipe left or right
func (s *Server) swipeHandler(w http.ResponseWriter, r *http.Request) {
	leftorright := chi.URLParam(r, "leftorright")

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
	if len(reqBody.Comment) > 256 {
		writeErrorResponse(w, r.Method, http.StatusBadRequest, "comment too long")
	}

	// left and right do the same thing for now
	// always return a response back to client, don't let them know about rabbitmq
	switch leftorright {
	case "left":
		writeStatusResponse(w, http.StatusCreated)
		s.IncrementThroughput()
	case "right":
		writeStatusResponse(w, http.StatusCreated)
		s.IncrementThroughput()
	default:
		writeErrorResponse(w, r.Method, http.StatusBadRequest, "not left or right")
		return
	}

	if err = s.publishToRmq(reqBody); err != nil {
		log.Error().Msgf("failed to publish to rabbitmq: %v", err)
	}
}

func (s *Server) publishToRmq(payload interface{}) error {
	log.Info().Interface("message", payload).Send()
	respBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return s.Pub.Publish(
		[]byte(respBytes),
		[]string{""},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange("swipes"),
	)
}

// Write a simple HTTP status to the response writer
func writeStatusResponse(w http.ResponseWriter, statusCode int) {
	log.Info().Int("code", statusCode)
	w.WriteHeader(statusCode)
}

// Marshal and write a JSON response to the response writer
func writeJsonResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	log.Info().Interface("send", payload).Send()
	respBytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "error marshaling JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(respBytes)))
	w.WriteHeader(statusCode)
	w.Write(respBytes)
}

// Write an HTTP error to the response writer
func writeErrorResponse(w http.ResponseWriter, method string, statusCode int, message string) {
	log.Error().Str("method", method).Int("code", statusCode).Msg(message)
	http.Error(w, message, statusCode)
}

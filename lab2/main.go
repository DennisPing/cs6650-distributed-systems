package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
)

func main() {
	log.Print("starting server...")

	// Create a new server
	server := NewServer()

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), server))
}

// Server implements chi.Router
type Server struct {
	store *InMemoryStore
	*chi.Mux
}

func NewServer() *Server {
	s := &Server{
		store: NewInMemoryStore(),
		Mux:   chi.NewRouter(),
	}
	s.Get("/", s.homeHandler)
	s.Route("/skiers", func(router chi.Router) {
		router.Get("/{skierID}/vertical", s.skierTotalVertHandler)
		router.Route("/{resortID}/seasons/{seasonID}/days/{dayID}/skiers/{skierID}", func(router chi.Router) {
			router.Get("/", s.skierDayVertHandler)
			router.Post("/", s.skierDayVertHandler)
		})
	})

	// For debugging purposes
	// chi.Walk(s, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
	// 	log.Printf("[%s]: '%s' has %d middlewares\n", method, route, len(middlewares))
	// 	return nil
	// })
	return s
}

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}

func (s *Server) skierTotalVertHandler(w http.ResponseWriter, r *http.Request) {
	skierIDTxt := chi.URLParam(r, "skierID")
	skierID, err := strconv.ParseInt(skierIDTxt, 10, 32)
	if err != nil {
		http.Error(w, "invalid skierID", http.StatusBadRequest)
		return
	}

	// Parse query params as pointers
	resortIDTxt := r.URL.Query().Get("resortID")
	var resortID *int32
	if resortIDTxt != "" {
		resortIDInt, err := strconv.ParseInt(resortIDTxt, 10, 32)
		if err != nil {
			http.Error(w, "invalid resortID", http.StatusBadRequest)
			return
		}
		resortIDValue := int32(resortIDInt)
		resortID = &resortIDValue
	}

	seasonIDTxt := r.URL.Query().Get("seasonID")
	var seasonID *string
	if seasonIDTxt != "" {
		seasonID = &seasonIDTxt
	}

	totalVert := s.store.GetSkierTotalVert(int32(skierID), resortID, seasonID)
	totalVertResp := TotalVertResponse{
		TotalVert: totalVert,
	}
	respJson, err := json.Marshal(totalVertResp)
	if err != nil {
		http.Error(w, "error marshaling JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respJson)
}

func (s *Server) skierDayVertHandler(w http.ResponseWriter, r *http.Request) {
	resortIDTxt := chi.URLParam(r, "resortID")
	seasonID := chi.URLParam(r, "seasonID")
	dayID := chi.URLParam(r, "dayID")
	skierIDTxt := chi.URLParam(r, "skierID")

	resortID, err := strconv.ParseInt(resortIDTxt, 10, 32)
	if err != nil {
		http.Error(w, "invalid resortID", http.StatusBadRequest)
		return
	}

	skierID, err := strconv.ParseInt(skierIDTxt, 10, 32)
	if err != nil {
		http.Error(w, "invalid skierID", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		totalVert := s.store.GetSkierDayVert(int32(resortID), seasonID, dayID, int32(skierID))
		totalVertResp := TotalVertResponse{
			TotalVert: totalVert,
		}
		respJson, err := json.Marshal(totalVertResp)
		if err != nil {
			http.Error(w, "error marshaling JSON response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respJson)
	} else if r.Method == http.MethodPost {
		var skier Skier
		decoder := json.NewDecoder(r.Body)
		// We don't do anything with the body for now...
		if err := decoder.Decode(&skier); err != nil {
			http.Error(w, "error decoding JSON request", http.StatusBadRequest)
			return
		}
		s.store.AddSkier(skier)
		w.WriteHeader(http.StatusOK)
	}
}

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
)

type Number interface {
	int32 | int64
}

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
	skierID := chi.URLParam(r, "skierID")
	resortID := r.URL.Query().Get("resortID") // optional
	seasonID := r.URL.Query().Get("seasonID") // optional

	skierIDNum, err := parseParamAsNumber[int32](skierID)
	if err != nil {
		writeError(w, r.Method, http.StatusBadRequest, err.Error())
		return
	}
	resortIDPtr, err := parseParamAsNumberPtr[int32](resortID)
	if err != nil {
		writeError(w, r.Method, http.StatusBadRequest, err.Error())
		return
	}
	seasonIDPtr, err := parseQueryParamAsStringPtr(seasonID)
	if err != nil {
		writeError(w, r.Method, http.StatusBadRequest, err.Error())
		return
	}

	totalVert := s.store.GetSkierTotalVert(skierIDNum, resortIDPtr, seasonIDPtr)
	totalVertResp := TotalVertResponse{
		TotalVert: totalVert,
	}
	writeJsonResponse(w, http.StatusOK, totalVertResp)
}

func (s *Server) skierDayVertHandler(w http.ResponseWriter, r *http.Request) {
	resortID := chi.URLParam(r, "resortID")
	seasonID := chi.URLParam(r, "seasonID")
	dayID := chi.URLParam(r, "dayID")
	skierID := chi.URLParam(r, "skierID")

	resortIDNum, err := parseParamAsNumber[int32](resortID)
	if err != nil {
		writeError(w, r.Method, http.StatusBadRequest, err.Error())
		return
	}
	skierIDNum, err := parseParamAsNumber[int32](skierID)
	if err != nil {
		writeError(w, r.Method, http.StatusBadRequest, err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		totalVert := s.store.GetSkierDayVert(resortIDNum, seasonID, dayID, skierIDNum)
		totalVertResp := TotalVertResponse{
			TotalVert: totalVert,
		}
		writeJsonResponse(w, http.StatusOK, totalVertResp)
	case http.MethodPost:
		var skier Skier
		decoder := json.NewDecoder(r.Body)
		// We don't do anything with the body for now...
		if err := decoder.Decode(&skier); err != nil {
			writeError(w, r.Method, http.StatusBadRequest, err.Error())
			return
		}
		s.store.AddSkier(skier)
		w.WriteHeader(http.StatusOK)
	default:
		writeError(w, r.Method, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// Marshal and write a JSON response to the response writer
func writeJsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
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
func writeError(w http.ResponseWriter, method string, statusCode int, message string) {
	log.Printf("Error: %s, %s, %d", method, message, statusCode)
	http.Error(w, message, statusCode)
}

// Generic function to parse a query param into a number
func parseParamAsNumber[T Number](param string) (T, error) {
	num, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return T(0), errors.New("not a number")
	}
	return T(num), nil
}

// Generic function to parse a query param into a number pointer
func parseParamAsNumberPtr[T Number](param string) (*T, error) {
	if param == "" {
		return nil, nil
	}
	num, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return nil, errors.New("not a number")
	}
	tNum := T(num)
	return &tNum, nil
}

// Parse a query param into a string pointer
func parseQueryParamAsStringPtr(param string) (*string, error) {
	if param == "" {
		return nil, nil
	}
	return &param, nil
}

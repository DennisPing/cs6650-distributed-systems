package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	log.Printf("starting server...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	http.HandleFunc("/", helloHandler)

	log.Printf("listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second) // Simulate work
	w.Write([]byte("Hello, world!"))
}

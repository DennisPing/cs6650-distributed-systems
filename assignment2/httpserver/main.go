package main

import (
	"fmt"
	"os"

	"github.com/DennisPing/cs6650-distributed-systems/assignment2/httpserver/log"
	"github.com/DennisPing/cs6650-distributed-systems/assignment2/httpserver/metrics"
	"github.com/DennisPing/cs6650-distributed-systems/assignment2/httpserver/rmq"
	"github.com/DennisPing/cs6650-distributed-systems/assignment2/httpserver/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)

	// Initialize metrics client
	metrics, err := metrics.NewMetrics()
	if err != nil {
		log.Fatal().Msgf("unable to set up metrics: %v", err)
	}

	// Initialize rabbitmq publisher
	rmqConn, err := rmq.NewConnection()
	if err != nil {
		log.Fatal().Msgf("unable to make rabbitmq connection: %v", err)
	}
	defer rmqConn.Close()
	pub, err := rmq.NewPublisher(rmqConn)
	if err != nil {
		log.Fatal().Msgf("unable to make rabbitmq publisher: %v", err)
	}
	defer pub.Close()

	// Start the http server
	server := server.NewServer(addr, metrics, pub)

	fmt.Printf("Starting server on port %s...", port)
	if err = server.Start(); err != nil {
		log.Fatal().Msgf("server died: %v", err)
	}
}

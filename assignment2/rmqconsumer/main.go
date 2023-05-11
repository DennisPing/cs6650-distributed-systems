package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DennisPing/cs6650-distributed-systems/assignment2/consumer/logger"
	"github.com/DennisPing/cs6650-distributed-systems/assignment2/consumer/rmq"
	"github.com/DennisPing/cs6650-distributed-systems/assignment2/consumer/server"
	"github.com/DennisPing/cs6650-distributed-systems/assignment2/consumer/store"
)

func main() {
	kvStore := store.NewStore()

	rmqConn, err := rmq.NewRmqConn()
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to make rabbitmq connection")
	}
	defer rmqConn.Close()

	consumer, err := rmq.StartRmqConsumer(rmqConn, kvStore)
	if err != nil {
		logger.Fatal().Err(err)
	}
	defer consumer.Close()

	// Start the HTTP server in the main goroutine, passing the KV store as a parameter
	server := server.Start(kvStore)

	// Set up a signal handler for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("shutting down gracefully...")
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		logger.Fatal().Err(err).Msg("failed to shutdown HTTP server gracefully")
	}
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/wagslane/go-rabbitmq"
)

var LogLevels = map[string]struct{}{
	"info":  {},
	"warn":  {},
	"error": {},
}

func main() {
	username := os.Getenv("RABBITMQ_USERNAME")
	password := os.Getenv("RABBITMQ_PASSWORD")
	host := os.Getenv("RABBITMQ_HOST")

	if username == "" || password == "" || host == "" {
		log.Fatal("You forgot to set the RABBITMQ env variables")
	}

	flag.Parse()
	if flag.NArg() == 0 {
		os.Exit(1)
	}
	var levels []string
	for _, level := range flag.Args() {
		if _, ok := LogLevels[level]; ok {
			levels = append(levels, level)
		} else {
			log.Fatalf("Invalid log level: %s", level)
		}
	}

	// Create a new connection to rabbitmq
	conn, err := rabbitmq.NewConn(
		fmt.Sprintf("amqp://%s:%s@%s:5672", username, password, host),
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	routingKeys := multipleRoutingKeys(levels)

	consumer, err := rabbitmq.NewConsumer(
		conn,
		func(d rabbitmq.Delivery) rabbitmq.Action {
			log.Printf("Received: %s", string(d.Body))
			return rabbitmq.Ack
		},
		"",
		append([]func(*rabbitmq.ConsumerOptions){
			rabbitmq.WithConsumerOptionsLogging,
			rabbitmq.WithConsumerOptionsExchangeDeclare,
			rabbitmq.WithConsumerOptionsExchangeDurable,
			rabbitmq.WithConsumerOptionsExchangeName("logs_direct"),
			rabbitmq.WithConsumerOptionsExchangeKind("direct"),
		}, routingKeys...)...,
	)
	if err != nil {
		log.Fatal(err)
	}

	defer consumer.Close()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down consumer gracefully")
}

func multipleRoutingKeys(keys []string) []func(*rabbitmq.ConsumerOptions) {
	options := make([]func(*rabbitmq.ConsumerOptions), 0, len(keys))
	for _, key := range keys {
		routingKeyOption := rabbitmq.WithConsumerOptionsRoutingKey(key)
		options = append(options, routingKeyOption)
	}
	return options
}

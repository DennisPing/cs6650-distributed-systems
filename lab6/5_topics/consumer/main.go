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
	var topics []string
	topics = append(topics, flag.Args()...)

	// Create a new connection to rabbitmq
	conn, err := rabbitmq.NewConn(
		fmt.Sprintf("amqp://%s:%s@%s:5672", username, password, host),
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	consumerOptions := createConsumerOptions(topics)

	consumer, err := rabbitmq.NewConsumer(
		conn,
		func(d rabbitmq.Delivery) rabbitmq.Action {
			log.Printf("Received: %s", string(d.Body))
			return rabbitmq.Ack
		},
		"",
		consumerOptions...,
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

// Create a slice of ptr ConsumerOptions using a slice of desired topics
func createConsumerOptions(topics []string) []func(*rabbitmq.ConsumerOptions) {
	baseOptions := []func(*rabbitmq.ConsumerOptions){
		rabbitmq.WithConsumerOptionsLogging,
		rabbitmq.WithConsumerOptionsExchangeDeclare,
		rabbitmq.WithConsumerOptionsExchangeDurable,
		rabbitmq.WithConsumerOptionsExchangeName("logs_topic"),
		rabbitmq.WithConsumerOptionsExchangeKind("topic"),
	}

	allOptions := append([]func(*rabbitmq.ConsumerOptions){}, baseOptions...)

	for _, topic := range topics {
		topicOption := rabbitmq.WithConsumerOptionsRoutingKey(topic)
		allOptions = append(allOptions, topicOption)
	}

	return allOptions
}

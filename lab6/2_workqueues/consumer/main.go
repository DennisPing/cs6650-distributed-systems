package main

import (
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

	// Create a new connection to rabbitmq
	conn, err := rabbitmq.NewConn(
		fmt.Sprintf("amqp://%s:%s@%s:5672", username, password, host),
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Create a new consumer that will consume messages from the queue
	consumer, err := rabbitmq.NewConsumer(
		conn,
		func(d rabbitmq.Delivery) rabbitmq.Action {
			log.Printf("Received: %s", string(d.Body))

			// Simulate work using time.Sleep()
			// dotCount := strings.Count(string(d.Body), ".")
			// time.Sleep(time.Duration(dotCount) * time.Second)
			return rabbitmq.Ack
		},
		"my_queue", // the mailbox to receive on
		rabbitmq.WithConsumerOptionsLogging,
		rabbitmq.WithConsumerOptionsQOSPrefetch(1),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Add a signal handler to gracefully shut down the consumer
	defer consumer.Close()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Closing consumer")
}

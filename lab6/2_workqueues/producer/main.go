package main

import (
	"bufio"
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

	conn, err := rabbitmq.NewConn(
		fmt.Sprintf("amqp://%s:%s@%s:5672", username, password, host),
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Create a new publisher with default settings
	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(""),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer publisher.Close()

	reader := bufio.NewReader(os.Stdin)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Start a persistent goroutine that publishes user messages
	go func() {
		for {
			fmt.Print("Enter message: ")
			msg, _ := reader.ReadString('\n')
			err := publisher.Publish(
				[]byte(msg),
				[]string{"my_queue"}, // the mailing address for this msg
				rabbitmq.WithPublishOptionsContentType("text/plain"),
				rabbitmq.WithPublishOptionsMandatory,
				rabbitmq.WithPublishOptionsExchange(""), // default exchange
			)
			if err != nil {
				log.Println(err)
			}
		}
	}()

	<-c
	fmt.Println("\nPublisher is shutting down gracefully")
}

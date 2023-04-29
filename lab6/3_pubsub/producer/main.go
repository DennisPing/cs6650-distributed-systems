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

	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName("logs"),
		rabbitmq.WithPublisherOptionsExchangeKind("fanout"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer publisher.Close()

	scanner := bufio.NewScanner(os.Stdin)

	go func() {
		fmt.Print(">> ")
		for scanner.Scan() {
			msg := scanner.Text()
			err := publisher.Publish(
				[]byte(msg),
				[]string{""},
				rabbitmq.WithPublishOptionsContentType("text/plain"),
				rabbitmq.WithPublishOptionsExchange("logs"),
			)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(">> ")
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Error reading from stdin: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down publisher gracefully")
}

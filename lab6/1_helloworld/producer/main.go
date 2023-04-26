package main

import (
	"fmt"
	"log"
	"os"

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
		rabbitmq.WithPublisherOptionsExchangeName("events"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer publisher.Close()

	err = publisher.Publish(
		[]byte("hello world!"),
		[]string{"my_routing_key"},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange("events"),
	)
	if err != nil {
		log.Println(err)
	}
}

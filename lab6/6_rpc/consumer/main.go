package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

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

	// Initialize the request channel
	reqChan := make(chan rabbitmq.Delivery)

	// Create the consumer who will receive the request
	consumer, err := rabbitmq.NewConsumer(
		conn,
		func(d rabbitmq.Delivery) rabbitmq.Action {
			reqChan <- d
			return rabbitmq.Ack
		},
		"rpc_queue",
		rabbitmq.WithConsumerOptionsLogging,
		rabbitmq.WithConsumerOptionsExchangeName(""),
		rabbitmq.WithConsumerOptionsExchangeKind("direct"),
		rabbitmq.WithConsumerOptionsQOSPrefetch(1),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	// Create the publisher who will send the response
	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(""),
		rabbitmq.WithPublisherOptionsExchangeKind("direct"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer publisher.Close()

	// Block until the request channel has a request to do
	for {
		d := <-reqChan
		var resp string
		n, err := strconv.Atoi(string(d.Body))
		if err != nil {
			resp = "Not a number"
		} else if n < 0 {
			resp = "Cannot be negative number"
		} else if n > 46 {
			resp = "Cannot be greater than 46"
		} else {
			fibNumber := fib(n)
			resp = strconv.Itoa(fibNumber)
			log.Printf("Computed fib(%d) = %d", n, fibNumber)
		}

		// Publish response back to the RPC caller
		err = publisher.Publish(
			[]byte(resp),
			[]string{d.ReplyTo},
			rabbitmq.WithPublishOptionsContentType("text/plain"),
			rabbitmq.WithPublishOptionsExchange(""),
			rabbitmq.WithPublishOptionsCorrelationID(d.CorrelationId),
		)
		if err != nil {
			log.Fatal(err)
		} else {
			log.Printf("Sending response: %s", resp)
		}
	}

}

// Inefficient recursive fibonacci formula
func fib(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fib(n-1) + fib(n-2)
	}
}

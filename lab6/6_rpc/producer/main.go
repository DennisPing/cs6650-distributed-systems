package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/wagslane/go-rabbitmq"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func main() {
	username := os.Getenv("RABBITMQ_USERNAME")
	password := os.Getenv("RABBITMQ_PASSWORD")
	host := os.Getenv("RABBITMQ_HOST")

	if username == "" || password == "" || host == "" {
		log.Fatal("You forgot to set the RABBITMQ env variables")
	}

	// Do input validation
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatal("Invalid args. Expected only 1 number.")
	}
	n, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		log.Fatal("Not a number")
	}

	// Create a new connection to rabbitmq
	conn, err := rabbitmq.NewConn(
		fmt.Sprintf("amqp://%s:%s@%s:5672", username, password, host),
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Send the RPC and wait for the response
	resp, err := FibonacciRPC(conn, n)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Got: %s", resp)
}

func FibonacciRPC(conn *rabbitmq.Conn, fibNumber int) (resp string, err error) {
	defer conn.Close()

	// Initialize response channel
	responseChan := make(chan string)

	// Create the publisher who will send the request
	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(""),
		rabbitmq.WithPublisherOptionsExchangeKind("direct"),
	)
	if err != nil {
		return "", err
	}
	defer publisher.Close()

	// Generate unique random IDs for correlation and reply-to queue
	rand.Seed(time.Now().UnixNano())
	correlationId := randString(32)
	replyToQueue := "reply_to_" + randString(16)

	// Create the consumer who will receive the response
	consumer, err := rabbitmq.NewConsumer(
		conn,
		func(d rabbitmq.Delivery) rabbitmq.Action {
			resp := string(d.Body)
			if d.CorrelationId == correlationId {
				responseChan <- resp // Send the response to the resp channel
			}
			return rabbitmq.Ack
		},
		replyToQueue,
		rabbitmq.WithConsumerOptionsLogging,
		rabbitmq.WithConsumerOptionsExchangeName(""),
		rabbitmq.WithConsumerOptionsExchangeKind("direct"),
		rabbitmq.WithConsumerOptionsQueueAutoDelete,
	)
	if err != nil {
		return "", err
	}
	defer consumer.Close()

	log.Printf("Sending: %d", fibNumber)
	// Send the RPC request
	err = publisher.Publish(
		[]byte(strconv.Itoa(fibNumber)),
		[]string{"rpc_queue"},
		rabbitmq.WithPublishOptionsContentType("text/plain"),
		rabbitmq.WithPublishOptionsExchange(""),
		rabbitmq.WithPublishOptionsCorrelationID(correlationId),
		rabbitmq.WithPublishOptionsReplyTo(replyToQueue),
	)
	if err != nil {
		return "", err
	}

	// Block until the response is received
	calculatedValue := <-responseChan
	return calculatedValue, nil
}

func randString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

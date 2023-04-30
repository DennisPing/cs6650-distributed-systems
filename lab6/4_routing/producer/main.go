package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
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
		rabbitmq.WithPublisherOptionsExchangeName("logs_direct"),
		rabbitmq.WithPublisherOptionsExchangeKind("direct"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer publisher.Close()

	scanner := bufio.NewScanner(os.Stdin)

	input := make(chan string)
	go func() {
		for scanner.Scan() {
			input <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Error reading from stdin: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	for {
		fmt.Print(">> ")
		select {
		case logText := <-input:
			logText = strings.TrimSpace(logText)
			if logText == "" {
				continue
			}
			parts := strings.Split(logText, " ")
			if _, ok := LogLevels[strings.ToLower(parts[0])]; ok {
				level := parts[0]
				msg := strings.Join(parts[1:], " ")
				err := publisher.Publish(
					[]byte(msg),
					[]string{level},
					rabbitmq.WithPublishOptionsContentType("text/plain"),
					rabbitmq.WithPublishOptionsExchange("logs_direct"),
				)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				log.Println("Invalid log level")
			}
		case <-quit:
			fmt.Println("Shutting down publisher gracefully")
			return
		}
	}
}

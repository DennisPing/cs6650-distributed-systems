package rmq

import (
	"errors"
	"fmt"
	"os"

	"github.com/wagslane/go-rabbitmq"
)

func NewConnection() (*rabbitmq.Conn, error) {
	username := os.Getenv("RABBITMQ_USERNAME")
	password := os.Getenv("RABBITMQ_PASSWORD")
	host := os.Getenv("RABBITMQ_HOST")

	if username == "" || password == "" || host == "" {
		return nil, errors.New("you forgot to set the RABBITMQ env variables")
	}

	conn, err := rabbitmq.NewConn(
		fmt.Sprintf("amqp://%s:%s@%s:5672", username, password, host),
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func NewPublisher(conn *rabbitmq.Conn) (*rabbitmq.Publisher, error) {
	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeDeclare,
		rabbitmq.WithPublisherOptionsExchangeName("swipes"),
		rabbitmq.WithPublisherOptionsExchangeKind("fanout"),
	)
	if err != nil {
		return nil, err
	}
	return publisher, nil
}

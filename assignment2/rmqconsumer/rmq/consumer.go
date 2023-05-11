package rmq

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/DennisPing/cs6650-distributed-systems/assignment2/consumer/logger"
	"github.com/DennisPing/cs6650-distributed-systems/assignment2/consumer/models"
	"github.com/DennisPing/cs6650-distributed-systems/assignment2/consumer/store"
	"github.com/wagslane/go-rabbitmq"
)

func NewRmqConn() (*rabbitmq.Conn, error) {
	username := os.Getenv("RABBITMQ_USERNAME")
	password := os.Getenv("RABBITMQ_PASSWORD")
	host := os.Getenv("RABBITMQ_HOST")

	if username == "" || password == "" || host == "" {
		logger.Fatal().Msg("you forgot to set the RABBITMQ env variables")
	}

	// Create a new connection to rabbitmq
	return rabbitmq.NewConn(
		fmt.Sprintf("amqp://%s:%s@%s:5672", username, password, host),
		rabbitmq.WithConnectionOptionsLogging,
	)
}

func StartRmqConsumer(conn *rabbitmq.Conn, kvStore *store.SimpleStore) (*rabbitmq.Consumer, error) {
	return rabbitmq.NewConsumer(
		conn,
		func(d rabbitmq.Delivery) rabbitmq.Action {
			logger.Info().Msg(string(d.Body))

			var reqBody models.SwipePayload
			err := json.Unmarshal(d.Body, &reqBody)
			if err != nil {
				logger.Error().Msgf("bad request: %v", err)
				return rabbitmq.NackDiscard
			}

			if err != nil {
				logger.Error().Msgf("invalid userId: %v", err)
				return rabbitmq.NackDiscard
			}

			kvStore.Add(reqBody.Swiper, reqBody.Swipee, reqBody.Direction)
			return rabbitmq.Ack
		},
		"",
		rabbitmq.WithConsumerOptionsLogging,
		rabbitmq.WithConsumerOptionsRoutingKey(""), // Bind this default queue to default routing key
		rabbitmq.WithConsumerOptionsExchangeName("swipes"),
		rabbitmq.WithConsumerOptionsExchangeKind("fanout"),
		rabbitmq.WithConsumerOptionsConcurrency(10),
		rabbitmq.WithConsumerOptionsQueueAutoDelete, // Auto delete the queue upon disconnect
	)
}

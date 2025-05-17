package queue

import (
	"fmt"
	"os"

	channel_pool "go-api-docker/internal/common/rabbitmq/channel_pool"
	publisher "go-api-docker/internal/common/rabbitmq/publisher"

	"github.com/joho/godotenv"
)

func CreateNewPublisher(channelPoolCount int, queueList []channel_pool.QueueDeclare) (*publisher.Publisher, error) {
	urlConnect, err := getUrlConnect()
	if err != nil {
		return nil, err
	}

	channelPool, err := channel_pool.NewChannelPool(urlConnect, channelPoolCount, queueList) // move to ENV "amqp://guest:guest@localhost:5672"
	if err != nil {
		return nil, err
	}

	var middlewares []publisher.Middleware
	publisherItem := publisher.NewPublisher(channelPool, middlewares) // thin about midllewear

	return publisherItem, nil
}

func getUrlConnect() (string, error) {
	if err := godotenv.Load(); err != nil {
		return "", err
	}

	var (
		host     = os.Getenv("RABBIT_MQ_HOST")
		port     = os.Getenv("RABBIT_MQ_PORT")
		user     = os.Getenv("RABBIT_MQ_USER")
		password = os.Getenv("RABBIT_MQ_PASSWORD")
	)

	urlConnect := fmt.Sprintf("amqp://%s:%s@%s:%s",
		user,
		password,
		host,
		port,
	)

	return urlConnect, nil
}

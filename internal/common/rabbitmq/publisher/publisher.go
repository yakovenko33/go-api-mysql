package rabbitmq

import (
	"log"
	"time"

	channel_pool "go-api-docker/internal/common/rabbitmq/channel_pool"

	"github.com/rabbitmq/amqp091-go"
)

// Middleware allows you to implement handlers before Publish (loggers, metrics, etc.)
type Middleware func(next PublisherFunc) PublisherFunc

// PublisherFunc defines the signature of the publish function
type PublisherFunc func(params PublisherParams) error

type Publisher struct {
	pool        *channel_pool.ChannelPool
	middlewares []Middleware
}

type PublisherParams struct {
	Exchange   string
	RoutingKey string
	RetryCount int
	RetryDelay int16
	Msg        amqp091.Publishing
}

func NewPublisher(pool *channel_pool.ChannelPool, middlewares []Middleware) *Publisher {
	return &Publisher{
		pool:        pool,
		middlewares: middlewares,
	}
}

func (p *Publisher) CloseChannelPool() {
	p.pool.Close()
}

func (p *Publisher) Publish(params PublisherParams) error {
	fn := p.publishWithRetry

	// we are running through middleware
	for i := len(p.middlewares) - 1; i >= 0; i-- {
		fn = p.middlewares[i](fn)
	}

	return fn(params)
}

func (p *Publisher) publishWithRetry(params PublisherParams) error {
	var err error

	for i := 0; i <= params.RetryCount; i++ {
		ch, chErr := p.pool.Get()
		if chErr != nil {
			err = chErr
			log.Printf("[retry %d/%d] get channel error: %v\n", i, params.RetryCount, chErr)
			time.Sleep(time.Duration(params.RetryDelay) * time.Second)
			continue
		}

		err = ch.Publish(params.Exchange, params.RoutingKey, false, false, params.Msg)
		p.pool.Put(ch)

		if err == nil {
			return nil
		}

		log.Printf("[retry %d/%d] publish error: %v\n", i, params.RetryCount, err) // need correct log
		time.Sleep(time.Duration(params.RetryDelay) * time.Second)
	}

	return err
}

/*pool, err := rabbitmq.NewPool("amqp://guest:guest@localhost:5672/", 5)
if err != nil {
	log.Fatal(err)
}
defer pool.Close()

loggerMiddleware := func(next rabbitmq.PublisherFunc) rabbitmq.PublisherFunc {
	return func(exchange, key string, msg amqp091.Publishing) error {
		log.Printf("📤 Publishing to [%s] key [%s]", exchange, key)
		return next(exchange, key, msg)
	}
}

middlewares := [1]Middleware{loggerMiddleware}
publisher := rabbitmq.NewPublisher(pool, 3, 2*time.Second, middlewares)

err = publisher.Publish(
	"my-exchange",
	"my.routing.key",
	amqp091.Publishing{
		ContentType: "text/plain",
		Body:        []byte("background task"),
	},
)
if err != nil {
	log.Println("❌ Failed to publish message:", err)
}*/

package rabbitmq

import (
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

// Middleware allows you to implement handlers before Publish (loggers, metrics, etc.)
type Middleware func(next PublisherFunc) PublisherFunc

// PublisherFunc defines the signature of the publish function
type PublisherFunc func(exchange, key string, msg amqp091.Publishing) error

type Publisher struct {
	pool        *ChannelPool
	retryCount  int
	retryDelay  time.Duration
	middlewares []Middleware
}

func NewPublisher(pool *ChannelPool, retryCount int, retryDelay time.Duration, middlewares ...Middleware) *Publisher {
	return &Publisher{
		pool:        pool,
		retryCount:  retryCount,
		retryDelay:  retryDelay,
		middlewares: middlewares,
	}
}

func (p *Publisher) Publish(exchange, key string, msg amqp091.Publishing) error {
	fn := p.publishWithRetry

	// we are running through middleware
	for i := len(p.middlewares) - 1; i >= 0; i-- {
		fn = p.middlewares[i](fn)
	}

	return fn(exchange, key, msg)
}

func (p *Publisher) publishWithRetry(exchange, key string, msg amqp091.Publishing) error {
	var err error

	for i := 0; i <= p.retryCount; i++ {
		ch, chErr := p.pool.Get()
		if chErr != nil {
			err = chErr
			log.Printf("🐛 [retry %d/%d] get channel error: %v\n", i, p.retryCount, chErr)
			time.Sleep(p.retryDelay)
			continue
		}

		err = ch.Publish(exchange, key, false, false, msg)
		p.pool.Put(ch)

		if err == nil {
			return nil
		}

		log.Printf("🔁 [retry %d/%d] publish error: %v\n", i, p.retryCount, err)
		time.Sleep(p.retryDelay * time.Duration(i+1)) // экспоненциальная задержка
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

publisher := rabbitmq.NewPublisher(pool, 3, 2*time.Second, loggerMiddleware)

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

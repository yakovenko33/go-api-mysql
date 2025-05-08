package rabbitmq

import (
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type HandlerFunc func(msg amqp091.Delivery) error

type Consumer struct {
	pool          *ChannelPool
	queue         string
	handler       HandlerFunc
	retryExchange string
	maxRetry      int
}

func NewConsumer(pool *ChannelPool, queue string, handler HandlerFunc, retryExchange string, maxRetry int) *Consumer {
	return &Consumer{
		pool:          pool,
		queue:         queue,
		handler:       handler,
		retryExchange: retryExchange,
		maxRetry:      maxRetry,
	}
}

func (c *Consumer) Start() error {
	ch, err := c.pool.Get()
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		c.queue,
		"",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			err := c.handler(msg)
			if err != nil {
				log.Println("❌ Handler error:", err)
				c.handleRetry(msg)
			} else {
				_ = msg.Ack(false)
			}
		}
	}()

	return nil
}

// retry через повторную публикацию (DLX + delay exchange)
func (c *Consumer) handleRetry(msg amqp091.Delivery) {
	attempts := getRetryCount(msg)

	if attempts >= c.maxRetry {
		log.Println("💀 Max retry reached, sending to DLX")
		_ = msg.Reject(false) // отправится в DLX, если настроен
		return
	}

	headers := msg.Headers
	if headers == nil {
		headers = amqp091.Table{}
	}
	headers["x-retry-count"] = attempts + 1

	// повторная публикация с задержкой (в DLX exchange)
	channel, err := c.pool.conn.Channel()
	if err != nil {
		log.Println("❌ Failed to open retry channel:", err)
		_ = msg.Nack(false, true)
		return
	}
	defer channel.Close()

	delay := time.Duration((attempts+1)*2) * time.Second

	err2 := channel.Publish(
		c.retryExchange,
		msg.RoutingKey,
		false,
		false,
		amqp091.Publishing{
			Headers:      headers,
			ContentType:  msg.ContentType,
			Body:         msg.Body,
			DeliveryMode: amqp091.Persistent,
			Expiration:   delay.String(), // TTL-based delay
		},
	)

	if err2 != nil {
		log.Println("❌ Retry publish failed:", err2)
		_ = msg.Nack(false, true)
		return
	}

	_ = msg.Ack(false)
}

func getRetryCount(msg amqp091.Delivery) int {
	val, ok := msg.Headers["x-retry-count"]
	if !ok {
		return 0
	}
	if v, ok := val.(int32); ok {
		return int(v)
	}
	if v, ok := val.(int64); ok {
		return int(v)
	}
	if v, ok := val.(int); ok {
		return v
	}
	return 0
}

/*
consumer := rabbitmq.NewConsumer(pool, "my-queue", func(msg amqp091.Delivery) error {
	log.Println("📥 Got message:", string(msg.Body))
	// simulate error
	if string(msg.Body) == "fail" {
		return fmt.Errorf("simulated failure")
	}
	return nil
}, "retry-exchange", 3)

if err := consumer.Start(); err != nil {
	log.Fatal(err)
}
*/

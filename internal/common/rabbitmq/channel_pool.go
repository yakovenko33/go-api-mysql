package rabbitmq

import (
	"errors"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ChannelPool struct {
	conn     *amqp.Connection
	pool     chan *amqp.Channel
	maxSize  int
	mu       sync.Mutex
	url      string
	isClosed bool
}

func NewChannelPool(url string, maxSize int) (*ChannelPool, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channelPool := &ChannelPool{
		conn:    conn,
		pool:    make(chan *amqp.Channel, maxSize),
		url:     url,
		maxSize: maxSize,
	}

	for i := 0; i < maxSize; i++ {
		ch, err := conn.Channel()
		if err != nil {
			return nil, err
		}
		channelPool.pool <- ch
	}

	go channelPool.watchReconnect()

	return channelPool, nil
}

func (p *ChannelPool) Get() (*amqp.Channel, error) {
	if p.isClosed {
		return nil, errors.New("pool closed")
	}

	ch := <-p.pool
	if ch.IsClosed() {
		log.Println("🔁 Канал закрыт, создаём новый")
		return p.conn.Channel()
	}

	return ch, nil
}

func (p *ChannelPool) Put(ch *amqp.Channel) {
	if ch == nil || ch.IsClosed() {
		return
	}

	select {
	case p.pool <- ch:
	default:
		_ = ch.Close()
	}
}

func (p *ChannelPool) Close() error {
	p.isClosed = true
	close(p.pool)
	for ch := range p.pool {
		_ = ch.Close()
	}
	return p.conn.Close()
}

func (p *ChannelPool) watchReconnect() {
	for {
		if p.conn.IsClosed() {
			log.Println("RabbitMQ connection lost. Reconnecting...") // need log to file
			for {
				conn, err := amqp.Dial(p.url)
				if err == nil {
					p.mu.Lock()
					p.conn = conn
					p.pool = make(chan *amqp.Channel, p.maxSize)
					for i := 0; i < p.maxSize; i++ {
						ch, err := conn.Channel()
						if err == nil {
							p.pool <- ch
						}
					}
					p.mu.Unlock()
					log.Println("✅ Reconnected to RabbitMQ")
					break
				}
				time.Sleep(5 * time.Second)
			}
		}
		time.Sleep(2 * time.Second)
	}
}

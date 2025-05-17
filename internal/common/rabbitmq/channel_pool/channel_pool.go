package rabbitmq

import (
	"errors"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ChannelPool struct {
	conn      *amqp.Connection
	pool      chan *amqp.Channel
	maxSize   int
	mu        sync.Mutex
	url       string
	isClosed  bool
	queueList []QueueDeclare
}

type QueueDeclare struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

func NewChannelPool(url string, maxSize int, queueList []QueueDeclare) (*ChannelPool, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channelPool := &ChannelPool{
		conn:      conn,
		pool:      make(chan *amqp.Channel, maxSize),
		url:       url,
		maxSize:   maxSize,
		queueList: queueList,
	}

	for i := 0; i < maxSize; i++ {
		ch, err := conn.Channel()
		if err != nil {
			closeChannelPool(channelPool)
			return nil, err
		}
		if err := queueDeclare(ch, queueList); err != nil {
			closeChannelPool(channelPool)
			return nil, err
		}
		channelPool.pool <- ch
	}

	go channelPool.watchReconnect()

	return channelPool, nil
}

func queueDeclare(ch *amqp.Channel, queueList []QueueDeclare) error {
	for _, queue := range queueList {
		_, err := ch.QueueDeclare(
			queue.Name,       // имя очереди
			queue.Durable,    // durable
			queue.AutoDelete, // delete when unused
			queue.Exclusive,  // exclusive
			queue.NoWait,     // no-wait
			queue.Args,       // arguments
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func closeChannelPool(channelPool *ChannelPool) error {
	channelPool.isClosed = true
	close(channelPool.pool)
	for ch := range channelPool.pool {
		_ = ch.Close()
	}
	return channelPool.conn.Close()
}

func (p *ChannelPool) Get() (*amqp.Channel, error) {
	if p.isClosed {
		return nil, errors.New("pool closed")
	}

	ch := <-p.pool
	if ch.IsClosed() {
		log.Println("🔁 Канал закрыт, создаём новый")
		channel, err := p.conn.Channel()
		if queueDeclare(channel, p.queueList); err != nil {
			return nil, err
		}
		return channel, err
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
	return closeChannelPool(p)
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
						if err != nil {
							closeChannelPool(p)
							return
						}
						if err := queueDeclare(ch, p.queueList); err != nil {
							closeChannelPool(p)
							return
						}
						p.pool <- ch
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

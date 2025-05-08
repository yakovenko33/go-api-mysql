package rabbitmq

import (
	"log"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type WatchdogConfig struct {
	AMQPUrl       string
	RetryInterval time.Duration
	OnReconnect   func(conn *amqp091.Connection)
	OnDisconnect  func(err error)
}

// Watchdog следит за соединением с RabbitMQ и восстанавливает его при падении.
type Watchdog struct {
	cfg    WatchdogConfig
	mu     sync.Mutex
	closed bool
}

// NewWatchdog создаёт watchdog.
func NewWatchdog(cfg WatchdogConfig) *Watchdog {
	return &Watchdog{cfg: cfg}
}

// Start запускает наблюдение за соединением.
func (w *Watchdog) Start() {
	go func() {
		for {
			conn, err := amqp091.Dial(w.cfg.AMQPUrl)
			if err != nil {
				log.Println("❌ Failed to connect to RabbitMQ:", err)
				time.Sleep(w.cfg.RetryInterval)
				continue
			}

			log.Println("✅ Connected to RabbitMQ")

			if w.cfg.OnReconnect != nil {
				w.cfg.OnReconnect(conn)
			}

			// Ждём, пока соединение не упадёт
			closeNotify := conn.NotifyClose(make(chan *amqp091.Error))
			err = <-closeNotify
			if err != nil {
				log.Println("⚠️ Connection closed:", err)
				if w.cfg.OnDisconnect != nil {
					w.cfg.OnDisconnect(err)
				}
			}

			// Проверка, не остановлен ли watchdog вручную
			w.mu.Lock()
			if w.closed {
				w.mu.Unlock()
				_ = conn.Close()
				return
			}
			w.mu.Unlock()

			log.Println("🔄 Reconnecting to RabbitMQ...")
			time.Sleep(w.cfg.RetryInterval)
		}
	}()
}

// Stop завершает работу watchdog.
func (w *Watchdog) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.closed = true
}

/*
func main() {
	watchdog := rabbitmq.NewWatchdog(rabbitmq.WatchdogConfig{
		AMQPUrl:       "amqp://guest:guest@localhost:5672/",
		RetryInterval: 5 * time.Second,
		OnReconnect: func(conn *amqp091.Connection) {
			ch, _ := conn.Channel()
			rabbitmq.DeclareTopology(ch)

			pool := rabbitmq.NewPoolFromConnection(conn) // ты должен сделать такой конструктор

			consumer := rabbitmq.NewConsumer(pool, "main-queue", handlerFunc, "retry-queue", 3)
			if err := consumer.Start(); err != nil {
				log.Println("❌ Failed to start consumer:", err)
			}
		},
		OnDisconnect: func(err error) {
			log.Println("💔 Connection lost:", err)
		},
	})

	watchdog.Start()

	select {} // блокируем main
}
*/

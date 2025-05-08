package rabbitmq

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func DeclareTopology(ch *amqp091.Channel) error {
	// 1. DLX очередь (куда попадут окончательно мёртвые сообщения)
	if _, err := ch.QueueDeclare("dead-letter-queue", true, false, false, false, nil); err != nil {
		return err
	}

	// 2. Retry очередь (сообщения с TTL вернутся в основную очередь)
	retryArgs := amqp091.Table{
		"x-dead-letter-exchange":    "",
		"x-dead-letter-routing-key": "main-queue",
		"x-message-ttl":             int32(10000), // 10 сек
	}
	if _, err := ch.QueueDeclare("retry-queue", true, false, false, false, retryArgs); err != nil {
		return err
	}

	// 3. Основная очередь, с DLX на retry-queue
	mainArgs := amqp091.Table{
		"x-dead-letter-exchange":    "",
		"x-dead-letter-routing-key": "retry-queue",
	}
	if _, err := ch.QueueDeclare("main-queue", true, false, false, false, mainArgs); err != nil {
		return err
	}

	log.Println("✅ RabbitMQ topology declared")
	return nil
}

/*
conn, _ := amqp091.Dial("amqp://guest:guest@localhost:5672/")
ch, _ := conn.Channel()
defer ch.Close()

if err := rabbitmq.DeclareTopology(ch); err != nil {
	log.Fatal("Failed to declare topology:", err)
}
*/

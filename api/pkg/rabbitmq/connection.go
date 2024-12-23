package rabbitmq

import (
	"api/pkg/config"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection

func Connect() {
	rsn := config.GetRSN()
	var err error

	for i := 0; i < 5; i++ {
		conn, err = amqp.Dial(rsn)
		if err == nil {
			break
		}
		log.Printf("Retrying RabbitMQ connection (%d/5): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	log.Println("RabbitMQ connection established")
}

func GetChannel() *amqp.Channel {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	return ch
}

func Close() {
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
			log.Println("failed to close RabbitMQ connection:", err)
		}
	}(conn)
}

func CloseChannel(ch *amqp.Channel) {
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			log.Println("failed to close RabbitMQ channel:", err)
		}
	}(ch)
}

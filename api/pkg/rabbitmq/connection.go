package rabbitmq

import (
	"api/pkg/config"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func CheckConnection() error {
	conn, err := Connect()
	defer Close(conn)
	if err != nil {
		return err
	}
	return nil
}

func Connect() (*amqp.Connection, error) {
	rsn := config.GetRSN()
	var (
		conn *amqp.Connection
		err  error
	)
	for i := 0; i < 5; i++ {
		conn, err = amqp.Dial(rsn)
		if err == nil {
			break
		}
		log.Printf("Retrying RabbitMQ connection (%d/5): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Printf("failed to connect to RabbitMQ: %v\n", err)
		return nil, err
	}
	log.Println("RabbitMQ connection established")
	return conn, err
}

func GetChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	return ch
}

func Close(conn *amqp.Connection) {
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

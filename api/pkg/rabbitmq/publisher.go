package rabbitmq

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct{}

func NewProducer() *Producer {
	return &Producer{}
}

func (producer *Producer) Publish(exchange, routingKey string, event interface{}) error {
	conn, err := Connect()
	if err != nil {
		return err
	}
	defer Close(conn)
	channel := GetChannel(conn)
	defer CloseChannel(channel)

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = channel.Publish(exchange, routingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		return err
	}
	err = channel.Confirm(false)
	if err != nil {
		return err
	}
	return nil
}

package rabbitmq

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	Channel *amqp.Channel
}

func NewProducer() *Producer {
	return &Producer{Channel: GetChannel()}
}

func (producer *Producer) Publish(exchange, routingKey string, event interface{}) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return producer.Channel.Publish(exchange, routingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

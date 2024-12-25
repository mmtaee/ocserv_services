package handler

import (
	"api/pkg/rabbitmq"
	"errors"
	"github.com/oklog/ulid/v2"
	"log"
)

type Event struct {
	ID     string      `json:"id"`
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

type Processor struct {
	producer *rabbitmq.Producer
}

func NewProcessor(producer *rabbitmq.Producer) *Processor {
	return &Processor{producer: producer}
}

func NewEvent(action string, data interface{}) *Event {
	return &Event{
		ID:     ulid.Make().String(),
		Action: action,
		Data:   data,
	}
}

func (p *Processor) ProcessEvent(event *Event, exchange string) error {
	if event.ID == "" || event.Action == "" || event.Data == nil {
		return errors.New("invalid event: missing required fields")
	}
	switch exchange {
	case "ocserv":
		return p.HandleOcservEvent(event)
	case "log":
		return p.HandleLogEvent(event)
	default:
		log.Printf("unknown exchange type: %s", exchange)
	}
	return nil
}

func (p *Processor) HandleOcservEvent(event *Event) error {
	return p.producer.Publish("ocserv", "ocserv.*", event)
}

func (p *Processor) HandleLogEvent(event *Event) error {
	return p.producer.Publish("log", "log.*", event)
}

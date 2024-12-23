package handler

import (
	"api/pkg/rabbitmq"
	"errors"
	"github.com/oklog/ulid/v2"
	"log"
	"strings"
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
		log.Printf("[WARN] unknown exchange type: %s", exchange)
	}
	return nil
}

func (p *Processor) HandleOcservEvent(event *Event) error {
	var routingKey string
	switch {
	case strings.HasPrefix(event.Action, "group"):
		routingKey = "ocserv.group.*"
	case strings.HasPrefix(event.Action, "user"):
		routingKey = "ocserv.user.*"
	default:
		return errors.New("invalid ocserv event")
	}
	return p.producer.Publish("ocserv", routingKey, event)
}

func (p *Processor) HandleLogEvent(event *Event) error {
	return p.producer.Publish("log", "log.*", event)
}

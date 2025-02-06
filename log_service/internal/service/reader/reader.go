package reader

import (
	"context"
	"log_service/internal/service/action"
	"log_service/internal/service/sse"
	"os"
)

type Reader interface {
	Start(chan string)
	Cancel() error
}

type Service struct {
	Reader Reader
	Action *action.Action
	SSE    *sse.Server
}

func SetReader() *Service {
	var re Reader
	ctx, cancel := context.WithCancel(context.Background())

	act := action.NewAction()
	sseServer := sse.NewSSEServer()

	if os.Getenv("DOCKERIZED") == "true" {
		re = &DockerReader{
			ctx:    ctx,
			cancel: cancel,
		}
	} else {
		re = &JournaldReader{
			ctx:    ctx,
			cancel: cancel,
		}
	}
	return &Service{
		Reader: re,
		Action: act,
		SSE:    sseServer,
	}
}

func (reader *Service) StartFetch() {
	messageChan := make(chan string, 100)

	go reader.Reader.Start(messageChan)
	go reader.Action.Start()
	go reader.SSE.Start()

	for message := range messageChan {
		reader.Action.Ch <- message
		reader.SSE.Broadcast(message)
	}

}

func (reader *Service) ShotDown() error {
	reader.SSE.Shutdown()
	if err := reader.Reader.Cancel(); err != nil {
		return err
	}
	if err := reader.Action.Cancel(); err != nil {
		return err
	}
	return nil
}

package reader

import (
	"context"
	"github.com/mmtaee/go-oc-utils/logger"
	"github.com/segmentio/kafka-go"
	"process/internal/activity"
	"process/internal/calculator"
)

type BrokerReader struct{}

var broker *kafka.Reader

func (b *BrokerReader) Start(act *activity.Activity, calc *calculator.Calculator) {
	ch := make(chan []byte)
	go consumer(ch)

	for msg := range ch {
		act.Ch <- string(msg)
		calc.Ch <- string(msg)
	}
}

func (b *BrokerReader) Cancel() error {
	err := broker.Close()
	if err != nil {
		return err
	}
	return nil
}

func consumer(ch chan<- []byte) {
	broker = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "ocserv",
		GroupID: "proc-ocserv",
	})
	for {
		msg, err := broker.ReadMessage(context.Background())
		if err != nil {
			logger.Logf(logger.ERROR, "failed to read message: %v", err)
			continue
		}
		ch <- msg.Value
	}
}

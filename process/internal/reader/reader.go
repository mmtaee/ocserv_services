package reader

import (
	"os"
	"process/internal/activity"
	"process/internal/calculator"
)

type Reader interface {
	Start(act *activity.Activity, calc *calculator.Calculator)
	Cancel() error
}

type Service struct {
	Re Reader
}

func NewReaderService() *Service {
	var reader Reader

	if os.Getenv("DOCKERIZED") == "true" {
		reader = &BrokerReader{}
	} else {
		reader = &SystemdReader{}
	}
	return &Service{Re: reader}
}

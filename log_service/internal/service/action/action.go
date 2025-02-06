package action

import (
	"context"
	"github.com/mmtaee/go-oc-utils/logger"
	"log_service/internal/repository/activity"
	"log_service/internal/repository/stats"
	"strings"
)

type Action struct {
	Ch     chan string
	ctx    context.Context
	cancel context.CancelFunc
}

func NewAction() *Action {
	ctx, cancel := context.WithCancel(context.Background())

	return &Action{
		Ch:     make(chan string, 100),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (action *Action) Start() {
	logger.Info("Log action service started")

	actionMap := map[string]func(string){
		"disconnected": func(text string) {
			go stats.Calculator(action.ctx, text)
			go activity.SetDisconnect(action.ctx, text)
		},
		"failed authentication": func(text string) { go activity.SetFailed(action.ctx, text) },
		"user logged in":        func(text string) { go activity.SetConnect(action.ctx, text) },
	}

	for {
		select {
		case msg := <-action.Ch:
			for keyword, act := range actionMap {
				if strings.Contains(msg, keyword) {
					go act(msg)
				}
			}
		case <-action.ctx.Done():
			return
		}
	}
}

func (action *Action) Cancel() error {
	action.cancel()
	return nil
}

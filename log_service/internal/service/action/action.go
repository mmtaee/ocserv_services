package action

import (
	"github.com/mmtaee/go-oc-utils/logger"
	"log_service/internal/repository/activity"
	"log_service/internal/repository/stats"
	"strings"
)

type Action struct {
	Ch chan string
}

func NewAction() *Action {
	return &Action{
		Ch: make(chan string, 100),
	}
}

func (action *Action) Start() {
	logger.Info("Log action service started")
	actionMap := map[string]func(string){
		"disconnected":          func(text string) { go stats.Calculator(text); go activity.SetDisconnect(text) },
		"failed authentication": func(text string) { go activity.SetFailed(text) },
		"user logged in":        func(text string) { go activity.SetConnect(text) },
	}

	for msg := range action.Ch {
		for keyword, act := range actionMap {
			if strings.Contains(msg, keyword) {
				go act(msg)
			}
		}
	}
}

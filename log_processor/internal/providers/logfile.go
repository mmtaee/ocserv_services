package providers

import (
	"context"
	"github.com/hpcloud/tail"
	"github.com/mmtaee/go-oc-utils/logger"
	"io"
	"service_log/internal/activity"
	"service_log/internal/stats"
	"strings"
)

func LogFile(c context.Context, logFilePath string) {
	streams, err := tail.TailFile(logFilePath, tail.Config{
		Follow:    true,
		MustExist: true,
		Poll:      true,
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: io.SeekEnd,
		},
		Logger: tail.DiscardingLogger,
	})
	if err != nil {
		logger.CriticalF("tail err: %v", err)
	}

	go func() {
		for line := range streams.Lines {
			if line.Err != nil {
				logger.Logf(logger.ERROR, "Error reading line: %v", line.Err)
				continue
			}

			actionMap := map[string]func(string){
				"disconnected":          func(text string) { go stats.Calculator(text); go activity.SetDisconnect(text) },
				"failed authentication": func(text string) { go activity.SetFailed(text) },
				"user logged in":        func(text string) { go activity.SetConnect(text) },
			}

			for keyword, action := range actionMap {
				if strings.Contains(line.Text, keyword) {
					go action(line.Text)
				}
			}
		}

		<-c.Done()

		if err = streams.Stop(); err != nil {
			logger.Logf(logger.ERROR, "Error stopping tail: %v", err)
		}
		return
	}()
}

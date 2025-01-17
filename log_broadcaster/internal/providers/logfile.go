package providers

import (
	"context"
	"github.com/hpcloud/tail"
	"github.com/mmtaee/go-oc-utils/logger"
	"io"
)

func LogFile(c context.Context, logFilePath string, ch chan string) {
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
		defer func() {
			if err = streams.Stop(); err != nil {
				logger.Logf(logger.ERROR, "Error stopping tail: %v", err)
			}
		}()

		for {
			select {
			case <-c.Done():
				logger.Logf(logger.WARNING, "Context canceled, stopping log file processing")
				return
			case line, ok := <-streams.Lines:
				if !ok {
					logger.Logf(logger.WARNING, "Log file stream closed unexpectedly")
					return
				}
				if line.Err != nil {
					logger.Logf(logger.ERROR, "Error reading line: %v", line.Err)
					continue
				}
				ch <- line.Text
			}
		}
	}()

}

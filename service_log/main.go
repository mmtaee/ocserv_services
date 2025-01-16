package main

import (
	"flag"
	"fmt"
	"github.com/hpcloud/tail"
	"service_log/internal/activity"

	// TODO: after complete develop remove this
	"github.com/joho/godotenv"
	"github.com/mmtaee/go-oc-utils/database"
	"github.com/mmtaee/go-oc-utils/logger"
	"io"
	"os"
	"os/signal"
	"service_log/internal/stats"
	"strings"
	"syscall"
	"time"
)

func main() {
	var (
		logFilePath string
		debug       bool
	)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	flag.StringVar(&logFilePath, "log-file", "", "Log file path")
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	flag.Parse()

	if logFilePath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if _, err := os.Stat(logFilePath); err != nil {
		logger.CriticalF("Failed to open log file: %v", err)
	}

	if debug {
		err := godotenv.Load()
		if err != nil {
			logger.Log(logger.CRITICAL, fmt.Sprintf("Error loading .env file: %v", err))
		}
	}

	cfg := &database.DBConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Name:     os.Getenv("POSTGRES_NAME"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	}
	database.Connect(cfg, debug)
	defer database.Close()

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
				"disconnected":          func(text string) { go stats.Calculator(text, line.Time); go activity.SetDisconnect(text) },
				"failed authentication": func(text string) { go activity.SetFailed(text) },
				"user logged in":        func(text string) { go activity.SetConnect(text) },
			}

			for keyword, action := range actionMap {
				if strings.Contains(line.Text, keyword) {
					go action(line.Text)
				}
			}
		}
	}()

	<-sigCh
	fmt.Println()
	logger.Log(logger.WARNING, "Shutting down service ...")
	time.Sleep(3 * time.Second)
}

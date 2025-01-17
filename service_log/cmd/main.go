package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/mmtaee/go-oc-utils/database"
	"github.com/mmtaee/go-oc-utils/logger"
	"os"
	"os/signal"
	"service_log/internal/providers"
	"syscall"
	"time"
)

func main() {
	var (
		logFile bool
		journal bool
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	flag.BoolVar(&logFile, "file", false, "Process log from file path")
	flag.BoolVar(&journal, "journal", false, "Process log from journalctl command")
	flag.Parse()

	if !logFile && journal {
		flag.PrintDefaults()
		os.Exit(1)
	}

	cfg := &database.DBConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Name:     os.Getenv("POSTGRES_NAME"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	}
	database.Connect(cfg, false)
	defer database.Close()

	if logFile {
		logFilePath := os.Getenv("LOG_FILE_PATH")
		if logFilePath == "" {
			logFilePath = "/var/log/ocserv/ocserv.log"
		}
		if _, err := os.Stat(logFilePath); err != nil {
			logger.CriticalF("Failed to open log file: %v", err)
		}
		providers.LogFile(ctx, logFilePath)
	} else {
		providers.Journal(ctx)
	}

	<-signalChan
	fmt.Println()
	logger.Log(logger.WARNING, "Shutting down service ...")
	cancel()
	time.Sleep(1 * time.Second)
}

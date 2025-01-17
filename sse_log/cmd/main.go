package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/mmtaee/go-oc-utils/logger"
	"net/http"
	"os"
	"os/signal"
	"sse_log/internal/events"
	"sse_log/internal/providers"
	"syscall"
	"time"
)

func main() {
	var (
		logFile bool
		journal bool
	)

	flag.BoolVar(&logFile, "file", false, "Process log from file path")
	flag.BoolVar(&journal, "journal", false, "Process log from journalctl command")
	flag.Parse()

	if !logFile && journal {
		flag.PrintDefaults()
		os.Exit(1)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	sse := events.NewSSEServer()
	defer sse.CloseAllClients()

	c, cancel := context.WithCancel(context.Background())
	defer cancel()

	if logFile {
		logFilePath := os.Getenv("LOG_FILE_PATH")
		if logFilePath == "" {
			logFilePath = "/var/log/ocserv/ocserv.log"
			logFilePath = "/tmp/ocserv/ocserv.log"
		}
		if _, err := os.Stat(logFilePath); err != nil {
			logger.CriticalF("Failed to open log file: %v", err)
		}
		providers.LogFile(c, logFilePath, sse.LogChan)
	} else {
		providers.Journal(c, sse.LogChan)
	}

	go func() {
		for msg := range sse.LogChan {
			sse.Broadcast(msg)
		}
	}()

	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: http.HandlerFunc(sse.ServerEventsHandler),
	}

	go func() {
		logger.InfoF("Starting server on %s:%s", host, port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.CriticalF("ListenAndServe error: %v", err)
		}
	}()

	<-signalChan
	logger.Log(logger.WARNING, "Shutdown signal received, cleaning up...")

	sse.CloseAllClients()

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.CriticalF("Server Shutdown error: %v", err)
	}
	time.Sleep(1 * time.Second)
	logger.Info("Shutdown complete")
}

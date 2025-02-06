package main

import (
	"fmt"
	"github.com/mmtaee/go-oc-utils/database"
	"github.com/mmtaee/go-oc-utils/logger"
	"log_service/internal/service/reader"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := &database.DBConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Name:     os.Getenv("POSTGRES_DB"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	}
	database.Connect(cfg, false)
	defer database.Close()

	logger.Info("Log service started")
	fetcher := reader.SetReader()
	go func() {
		fetcher.StartFetch()
	}()

	defer func(reader *reader.Service) {
		err := reader.ShotDown()
		if err != nil {
			logger.CriticalF("Error during cancel: %v", err)
		}
	}(fetcher)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	sig := <-quit
	fmt.Println()
	logger.InfoF("Signal received: %s. Initiating shutdown...", sig)

	err := fetcher.ShotDown()
	if err != nil {
		logger.CriticalF("Error during cancel: %v", err)
	}

	logger.Info("Shutdown complete")
}

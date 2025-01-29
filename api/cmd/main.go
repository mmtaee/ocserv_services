package main

import (
	_ "api/docs"
	"api/internal/handlers"
	"api/pkg/config"
	"api/pkg/event"
	"api/pkg/routing"
	"flag"
	"fmt"
	"github.com/mmtaee/go-oc-utils/database"
	"github.com/mmtaee/go-oc-utils/logger"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

// @title Ocserv User management Example Api
// @version 1.0
// @description This is a sample Ocserv User management Api server.
// @BasePath /api
func main() {
	var (
		debug   bool
		migrate bool
		drop    bool
	)
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.BoolVar(&migrate, "migrate", false, "migrate models to database")
	flag.BoolVar(&drop, "drop", false, "drop models table from database")
	flag.Parse()

	config.Set(debug)
	dbCfg := config.GetDB()
	dbConfig := &database.DBConfig{
		Host:     dbCfg.Host,
		Port:     dbCfg.Port,
		User:     dbCfg.User,
		Password: dbCfg.Password,
		Name:     dbCfg.Name,
	}
	database.Connect(dbConfig, debug)
	if migrate {
		handlers.Migrate()
	} else if drop && debug {
		handlers.Drop()
	} else {
		var eventWorkerCount int
		if eventWorkerStr := os.Getenv("EVENT_WORKER"); eventWorkerStr != "" {
			var err error
			eventWorkerCount, err = strconv.Atoi(eventWorkerStr)
			if err != nil {
				eventWorkerCount = 1
			}
		} else {
			eventWorkerCount = 1
		}

		go func() {
			routing.Serve()
		}()

		event.Set(database.Connection(), 100)
		eventWorker := event.GetWorker()
		go func() {
			eventWorker.Start(eventWorkerCount)
		}()

		defer func() {
			eventWorker.Stop()
			routing.Shutdown()
			database.Close()
			logger.Info("Shutdown complete")
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

		sig := <-quit
		fmt.Println()
		logger.InfoF("Signal received: %s. Initiating shutdown...", sig)

		eventWorker.Stop()
		routing.Shutdown()
		database.Close()
		logger.Info("Shutdown complete")
	}
}

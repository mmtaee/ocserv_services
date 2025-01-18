package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/mmtaee/go-oc-utils/database"
	"github.com/mmtaee/go-oc-utils/logger"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"user_expiry/checker"
)

var (
	restore bool
	expire  bool
)

func main() {
	flag.BoolVar(&restore, "restore", false, "Restore expired user")
	flag.BoolVar(&expire, "expire", false, "Expire user account")
	flag.Parse()

	if !restore && !expire {
		logger.Log(logger.ERROR, "one of -restore, -expire or -restore must be set")
		flag.PrintDefaults()
		os.Exit(0)
	}

	cfg := &database.DBConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Name:     os.Getenv("POSTGRES_DB"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	}
	database.Connect(cfg, false)
	defer database.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	if restore {
		wg.Add(1)
		go func() {
			defer wg.Done()
			checker.RestoreMonthlyAccounts(ctx)
		}()
	}
	if expire {
		wg.Add(1)
		go func() {
			defer wg.Done()
			checker.CheckExpiry(ctx)
		}()
	}

	<-signalChan
	fmt.Println()
	logger.Log(logger.WARNING, "Shutting down service ...")

	cancel()
	wg.Wait()

	time.Sleep(1 * time.Second)
	logger.Info("Shutdown complete")
}

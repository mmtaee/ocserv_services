package main

import (
	"fmt"
	"github.com/mmtaee/go-oc-utils/logger"
	"os"
	"os/signal"
	"process/internal/activity"
	"process/internal/calculator"
	"process/internal/reader"
	"syscall"
)

func main() {
	re := reader.NewReaderService()
	act := activity.NewActivityService()
	calc := calculator.NewCalculator()

	go func() {
		re.Re.Start(act, calc)
	}()

	defer func(Re reader.Reader) {
		err := Re.Cancel()
		if err != nil {
			logger.CriticalF("Error during cancel: %v", err)
		}
	}(re.Re)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	sig := <-quit
	fmt.Println()
	logger.InfoF("Signal received: %s. Initiating shutdown...", sig)

	err := re.Re.Cancel()
	if err != nil {
		logger.CriticalF("Error during cancel: %v", err)
	}

	logger.Info("Shutdown complete")
}

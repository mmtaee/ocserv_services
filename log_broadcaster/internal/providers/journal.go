package providers

import (
	"bufio"
	"context"
	"github.com/mmtaee/go-oc-utils/logger"
	"os/exec"
)

func Journal(c context.Context, ch chan string) {
	cmd := exec.Command("journalctl", "-fu", "ocserv")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.CriticalF("Error getting stdout pipe: %v\n", err)
	}

	if err = cmd.Start(); err != nil {
		logger.CriticalF("Error starting journalctl command: %v\n", err)
		return
	}
	scanner := bufio.NewScanner(stdout)
	logger.Info("Streaming logs for service: ocserv")

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			if err = scanner.Err(); err != nil {
				logger.CriticalF("Error reading journalctl output: %v", err)
			}
			ch <- line
		}
		<-c.Done()

		if err = cmd.Process.Kill(); err != nil {
			logger.CriticalF("Error killing journalctl process: %v", err)
		}

		if err = cmd.Wait(); err != nil {
			logger.CriticalF("journalctl command exited with error: %v", err)
		}
	}()
}

package reader

import (
	"bufio"
	"context"
	"github.com/mmtaee/go-oc-utils/logger"
	"os/exec"
)

type JournaldReader struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (j *JournaldReader) Start(ch chan string) {
	checkService := exec.Command("systemctl", "status", "ocserv.service")
	if err := checkService.Run(); err != nil {
		close(ch)
		logger.Critical("Ocserv service is not running")
	}

	cmd := exec.Command("journalctl", "-fq", "ocserv.service")
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
			if err = scanner.Err(); err != nil {
				logger.CriticalF("Error reading journalctl output: %v", err)
			}
			line := scanner.Text()
			ch <- line
		}
		<-j.ctx.Done()
		if err = cmd.Process.Kill(); err != nil {
			logger.CriticalF("Error killing journalctl process: %v", err)
		}
		if err = cmd.Wait(); err != nil {
			logger.CriticalF("journalctl command exited with error: %v", err)
		}
	}()
}

func (j *JournaldReader) Cancel() error {
	j.cancel()
	return nil
}

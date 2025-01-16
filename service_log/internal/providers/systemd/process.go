package systemd

import (
	"bufio"
	"context"
	"github.com/mmtaee/go-oc-utils/logger"
	"os/exec"
	"service_log/internal/activity"
	"service_log/internal/stats"
	"strings"
)

func Process(c context.Context) {
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
			actionMap := map[string]func(string){
				"disconnected":          func(text string) { go stats.Calculator(text); go activity.SetDisconnect(text) },
				"failed authentication": func(text string) { go activity.SetFailed(text) },
				"user logged in":        func(text string) { go activity.SetConnect(text) },
			}

			for keyword, action := range actionMap {
				if strings.Contains(line, keyword) {
					go action(line)
				}
			}
			if err = scanner.Err(); err != nil {
				logger.CriticalF("Error reading journalctl output: %v", err)
			}
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

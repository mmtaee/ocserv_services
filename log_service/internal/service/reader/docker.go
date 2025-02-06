package reader

import (
	"bufio"
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/mmtaee/go-oc-utils/logger"
	"regexp"
	"strings"
	"time"
)

type DockerReader struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (d *DockerReader) Start(ch chan string) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.CriticalF("Error creating Docker client:", err)
	}

	reader, err := cli.ContainerLogs(d.ctx, "ocserv-api", container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: false,
		Since:      time.Now().Format(time.RFC3339),
	})
	if err != nil {
		logger.CriticalF("Error getting logs:", err)
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "ocserv[") {
			msg := strings.ReplaceAll(scanner.Text(), "\x00", "")
			re := regexp.MustCompile(`^.*?ocserv\[\d+\]:\s*`)
			msg = re.ReplaceAllString(msg, "")
			ch <- msg
		}
	}
	if err = scanner.Err(); err != nil {
		logger.CriticalF("Error reading logs:", err)
	}
}

func (d *DockerReader) Cancel() error {
	d.cancel()
	return nil
}

package ocserv

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Handler struct {
	Group OcGroupInterface
	User  OcUserInterface
	Occtl OcctlInterface
}

func NewHandler() *Handler {
	return &Handler{
		Group: NewOcGroup(),
		User:  NewOcUser(),
		Occtl: NewOcctl(),
	}
}

var (
	ocpasswdCMD  = "sudo /usr/bin/ocpasswd"
	passwdFile   = "/etc/ocserv/ocpasswd"
	groupDir     = "/etc/ocserv/groups"
	defaultGroup = "/etc/ocserv/defaults/group.conf"
)

func WithContext(ctx context.Context, operation func() error) error {
	done := make(chan error, 1)

	go func() {
		defer close(done)
		done <- operation()
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("operation canceled or timed out: %w", ctx.Err())
	case err := <-done:
		return err
	}
}

func (h *Handler) ToMap(data interface{}) map[string]interface{} {
	b, _ := json.Marshal(&data)
	var dataStruct map[string]interface{}
	_ = json.Unmarshal(b, &dataStruct)
	return dataStruct
}

func (h *Handler) ReadOcpasswd() (userList []Sync) {
	content, err := os.ReadFile(passwdFile)
	if err != nil {
		return nil
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		userSplit := strings.Split(line, ":")
		if len(userSplit) == 2 {
			userList = append(userList, Sync{
				Username: userSplit[0],
				Group:    userSplit[1],
			})
		}
	}
	return userList
}
func ParseConfFile(filename string) (OcGroupConfig, error) {
	var config OcGroupConfig
	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "rx-data-per-sec":
			config.RxDataPerSec = value
		case "tx-data-per-sec":
			config.TxDataPerSec = value
		case "max-same-clients":
			if val, err := strconv.Atoi(value); err == nil {
				config.MaxSameClients = val
			}
		case "ipv4-network":
			config.IPv4Network = value
		case "dns":
			config.DNS = strings.Split(value, ",")
		case "no-udp":
			if val, err := strconv.ParseBool(value); err == nil {
				config.NoUDP = val
			}
		case "keepalive":
			if val, err := strconv.Atoi(value); err == nil {
				config.KeepAlive = val
			}
		case "dpd":
			if val, err := strconv.Atoi(value); err == nil {
				config.DPD = val
			}
		case "mobile-dpd":
			if val, err := strconv.Atoi(value); err == nil {
				config.MobileDPD = val
			}
		case "tunnel-all-dns":
			if val, err := strconv.ParseBool(value); err == nil {
				config.TunnelAllDNS = val
			}
		case "restrict-user-to-routes":
			if val, err := strconv.ParseBool(value); err == nil {
				config.RestrictUserToRoutes = val
			}
		case "stats-report-time":
			if val, err := strconv.Atoi(value); err == nil {
				config.StatsReportTime = val
			}
		case "mtu":
			if val, err := strconv.Atoi(value); err == nil {
				config.MTU = val
			}
		case "idle-timeout":
			if val, err := strconv.Atoi(value); err == nil {
				config.IdleTimeout = val
			}
		case "mobile-idle-timeout":
			if val, err := strconv.Atoi(value); err == nil {
				config.MobileIdleTimeout = val
			}
		case "session-timeout":
			if val, err := strconv.Atoi(value); err == nil {
				config.SessionTimeout = val
			}
		}
	}

	if err = scanner.Err(); err != nil {
		return config, err
	}

	return config, nil
}

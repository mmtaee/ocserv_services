package ocserv

import (
	"encoding/json"
	"log"
	"os/exec"
)

type Occtl struct {
}

type OcctlInterface interface {
	Reload() error
	OnlineUsers() []OcctlUser
	Disconnect(string) error
	ShowIPBans(bool) []IPBan
	UnBanIP(string) error
	ShowStatus() string
	ShowIRoutes() []IRoute
}

var occtlCMD = "sudo /usr/bin/occtl"

func NewOcctl() *Occtl {
	return &Occtl{}
}

func (o *Occtl) Reload() error {
	cmd := exec.Command(occtlCMD, "reload")
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	log.Println("Command Output:\n", string(output))
	return nil
}

func (o *Occtl) OnlineUsers() []OcctlUser {
	cmd := exec.Command(occtlCMD, "-j", "show", "users", "--output=json-pretty")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}
	var users []OcctlUser
	err = json.Unmarshal(output, &users)
	if err != nil {
		return nil
	}
	return users
}

func (o *Occtl) Disconnect(username string) error {
	cmd := exec.Command(occtlCMD, "disconnect", "user", username)
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	log.Println("Command Output:\n", string(output))
	return nil
}

func (o *Occtl) ShowIPBans(points bool) []IPBan {
	command := []string{"-j", "show", "ip", "bans"}
	if points {
		command = append(command, "points")
	}
	cmd := exec.Command(occtlCMD, command...)
	output, err := cmd.Output()
	if err != nil {
		return nil
	}
	var ipBans []IPBan
	err = json.Unmarshal(output, &ipBans)
	if err != nil {
		return nil
	}
	return ipBans
}

func (o *Occtl) UnBanIP(IP string) error {
	cmd := exec.Command(occtlCMD, "unban", "ip", IP)
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	log.Println("Command Output:\n", string(output))
	return nil
}

func (o *Occtl) ShowStatus() string {
	cmd := exec.Command(occtlCMD, "show", "status")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(output)
}

func (o *Occtl) ShowIRoutes() []IRoute {
	cmd := exec.Command(occtlCMD, "-j", "show", "iroutes")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}
	var routes []IRoute
	err = json.Unmarshal(output, &routes)
	if err != nil {
		return nil
	}
	return routes
}

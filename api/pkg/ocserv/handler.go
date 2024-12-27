package ocserv

import (
	"encoding/json"
	"os"
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

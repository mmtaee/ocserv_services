package ocserv

import (
	"fmt"
	"log"
	"os/exec"
)

type OcUser struct{}

type OcUserInterface interface {
	CreateOrUpdateUser(string, string, string) error
	DeleteUser(string) error
	LockUnLockUser(string, bool) error
}

func NewOcUser() *OcUser {
	return &OcUser{}
}

func (o *OcUser) CreateOrUpdateUser(username, password, group string) error {
	if group == "defaults" {
		group = ""
	} else {
		group = fmt.Sprintf("-g %s", group)
	}
	command := fmt.Sprintf("/usr/bin/echo -e \"%s\\n%s\\n | %s %s -c %s %s",
		password,
		password,
		ocpasswdCMD,
		group,
		passwdFile,
		username,
	)
	cmd := exec.Command(command)
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	log.Println("Command Output:\n", string(output))
	return nil
}

func (o *OcUser) DeleteUser(username string) error {
	command := fmt.Sprintf("%s -c %s -d %s", ocpasswdCMD, passwdFile, username)
	cmd := exec.Command(command)
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	log.Println("Command Output:\n", string(output))
	return nil
}

func (o *OcUser) LockUnLockUser(username string, lock bool) error {
	lockAction := "-u"
	if lock {
		lockAction = "-l"
	}
	command := fmt.Sprintf("%s %s -c %s %s", ocpasswdCMD, lockAction, passwdFile, username)
	cmd := exec.Command(command)
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	log.Println("Command Output:\n", string(output))
	return nil
}

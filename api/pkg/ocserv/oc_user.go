package ocserv

import (
	"context"
	"fmt"
	"log"
	"os/exec"
)

type OcUser struct{}

type OcUserInterface interface {
	CreateOrUpdateUser(c context.Context, username, password, group string) error
	LockUnLockUser(c context.Context, username string, lock bool) error
	DeleteUser(username string) error
}

func NewOcUser() *OcUser {
	return &OcUser{}
}

func (o *OcUser) CreateOrUpdateUser(c context.Context, username, password, group string) error {
	if group == "defaults" {
		group = ""
	} else {
		group = fmt.Sprintf("-g %s", group)
	}
	return WithContext(c, func() error {
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
	})
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

func (o *OcUser) LockUnLockUser(c context.Context, username string, lock bool) error {
	lockAction := "-u"
	if lock {
		lockAction = "-l"
	}
	return WithContext(c, func() error {
		command := fmt.Sprintf("%s %s -c %s %s", ocpasswdCMD, lockAction, passwdFile, username)
		cmd := exec.Command(command)
		output, err := cmd.Output()
		if err != nil {
			return err
		}
		log.Println("Command Output:\n", string(output))
		return nil
	})
}

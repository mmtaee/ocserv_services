package ocserv

import (
	"context"
	"fmt"
	"os/exec"
)

type OcUser struct{}

type OcUserInterface interface {
	CreateOrUpdateUser(c context.Context, username, password, group string) error
	LockUnLockUser(c context.Context, username string, lock bool) error
	DeleteUser(c context.Context, username string) error
}

func NewOcUser() *OcUser {
	return &OcUser{}
}

func (o *OcUser) CreateOrUpdateUser(c context.Context, username, password, group string) error {
	if group == "defaults" || group == "" {
		group = ""
	} else {
		group = fmt.Sprintf("-g %s", group)
	}
	command := fmt.Sprintf("/usr/bin/echo -e \"%s\\n%s\\n\" | %s %s -c %s %s",
		password,
		password,
		ocpasswdCMD,
		group,
		passwdFile,
		username,
	)
	return exec.CommandContext(c, "sh", "-c", command).Run()
}

func (o *OcUser) LockUnLockUser(c context.Context, username string, lock bool) error {
	lockAction := "-u"
	if lock {
		lockAction = "-l"
	}
	command := fmt.Sprintf("%s %s -c %s %s", ocpasswdCMD, lockAction, passwdFile, username)
	return exec.CommandContext(c, "sh", "-c", command).Run()
}

func (o *OcUser) DeleteUser(c context.Context, username string) error {
	command := fmt.Sprintf("%s -c %s -d %s", ocpasswdCMD, passwdFile, username)
	return exec.CommandContext(c, "sh", "-c", command).Run()
}

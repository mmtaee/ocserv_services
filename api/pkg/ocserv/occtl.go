package ocserv

import (
	"context"
	"encoding/json"
	"fmt"
)

type Occtl struct {
}

type OcctlInterface interface {
	Reload(c context.Context) error
	OnlineUsers(c context.Context) (*[]OcctlUser, error)
	Disconnect(c context.Context, username string) error
	ShowIPBans(c context.Context, points bool) []IPBan
	UnBanIP(c context.Context, ip string) error
	ShowStatus(c context.Context) string
	ShowIRoutes(c context.Context) []IRoute
	ShowUser(c context.Context, username string) (*OcctlUser, error)
}

func NewOcctl() *Occtl {
	return &Occtl{}
}

func (o *Occtl) Reload(c context.Context) error {
	_, err := OcctlExec(c, "reload")
	if err != nil {
		return err
	}
	return nil
}

func (o *Occtl) OnlineUsers(c context.Context) (*[]OcctlUser, error) {
	var users []OcctlUser
	result, err := OcctlExec(c, "-j show users --output=json-pretty")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(result, &users)
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func (o *Occtl) Disconnect(c context.Context, username string) error {
	_, err := OcctlExec(c, fmt.Sprintf("disconnect user %s", username))
	if err != nil {
		return err
	}
	return nil
}

func (o *Occtl) ShowIPBans(c context.Context, points bool) []IPBan {
	command := "-j show ip bans"
	if points {
		command += " points"
	}
	result, err := OcctlExec(c, command)
	if err != nil {
		return nil
	}
	var ipBans []IPBan
	err = json.Unmarshal(result, &ipBans)
	if err != nil {
		return nil
	}
	return ipBans
}

func (o *Occtl) UnBanIP(c context.Context, ip string) error {
	_, err := OcctlExec(c, fmt.Sprintf("unban ip %s", ip))
	if err != nil {
		return err
	}
	return nil
}

func (o *Occtl) ShowStatus(c context.Context) string {
	result, err := OcctlExec(c, "show status")
	if err != nil {
		return ""
	}
	return string(result)
}

func (o *Occtl) ShowIRoutes(c context.Context) []IRoute {
	result, err := OcctlExec(c, "-j show iroutes")
	if err != nil {
		return nil
	}
	var routes []IRoute
	err = json.Unmarshal(result, &routes)
	if err != nil {
		return nil
	}
	return routes
}

func (o *Occtl) ShowUser(c context.Context, username string) (*OcctlUser, error) {
	var user *OcctlUser
	result, err := OcctlExec(c, fmt.Sprintf("-j show user %s", username))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(result, &user)
	if err != nil {
		return nil, err
	}
	return user, err
}

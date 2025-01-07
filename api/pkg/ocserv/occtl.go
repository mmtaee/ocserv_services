package ocserv

import (
	"context"
	"encoding/json"
	"fmt"
)

type Occtl struct {
}

type OcctlInterface interface {
	Reload() error
	OnlineUsers(c context.Context) ([]OcctlUser, error)
	Disconnect(c context.Context, username string) error
	ShowIPBans(bool) []IPBan
	UnBanIP(string) error
	ShowStatus() string
	ShowIRoutes() []IRoute
	ShowUser(c context.Context, username string) (OcctlUser, error)
}

func NewOcctl() *Occtl {
	return &Occtl{}
}

func (o *Occtl) Reload() error {
	_, err := OcctlExec("reload")
	if err != nil {
		return err
	}
	return nil
}

func (o *Occtl) OnlineUsers(c context.Context) ([]OcctlUser, error) {
	var users []OcctlUser
	err := WithContext(c, func() error {
		result, err := OcctlExec("-j show users --output=json-pretty")
		if err != nil {
			return err
		}
		return json.Unmarshal(result, &users)
	})
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (o *Occtl) Disconnect(c context.Context, username string) error {
	return WithContext(c, func() error {
		_, err := OcctlExec(fmt.Sprintf("disconnect user %s", username))
		if err != nil {
			return err
		}
		return nil
	})
}

func (o *Occtl) ShowIPBans(points bool) []IPBan {
	command := "-j show ip bans"
	if points {
		command += " points"
	}
	result, err := OcctlExec(command)
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

func (o *Occtl) UnBanIP(ip string) error {
	_, err := OcctlExec(fmt.Sprintf("unban ip %s", ip))
	if err != nil {
		return err
	}
	return nil
}

func (o *Occtl) ShowStatus() string {
	result, err := OcctlExec("show status")
	if err != nil {
		return ""
	}
	return string(result)
}

func (o *Occtl) ShowIRoutes() []IRoute {
	result, err := OcctlExec("-j show iroutes")
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

func (o *Occtl) ShowUser(c context.Context, username string) (OcctlUser, error) {
	var user OcctlUser
	err := WithContext(c, func() error {
		result, err := OcctlExec(fmt.Sprintf("-j show user %s", username))
		if err != nil {
			return err
		}
		return json.Unmarshal(result, &user)
	})
	return user, err
}

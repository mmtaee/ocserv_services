package repository

import (
	"context"
	"github.com/mmtaee/go-oc-utils/handler/occtl"
)

type OcctlRepository struct {
	oc occtl.OcInterface
}

type OcctlRepositoryInterface interface {
	Reload(c context.Context) error
	OnlineUsers(c context.Context) (*[]occtl.OcUser, error)
	Disconnect(c context.Context, username string) error
	ShowIPBans(c context.Context) (*[]occtl.IPBan, error)
	ShowIPBansPoint(c context.Context) (*[]occtl.IPBanPoints, error)
	UnBanIP(c context.Context, ip string) error
	ShowStatus(c context.Context) string
	ShowIRoutes(c context.Context) (*[]occtl.IRoute, error)
	ShowUser(c context.Context, username string) (*[]occtl.OcUser, error)
}

func NewOcctlRepository() *OcctlRepository {
	return &OcctlRepository{oc: occtl.NewOcctl()}
}

func (o *OcctlRepository) Reload(c context.Context) error {
	return o.oc.Reload(c)
}

func (o *OcctlRepository) OnlineUsers(c context.Context) (*[]occtl.OcUser, error) {
	return o.oc.OnlineUsers(c)
}

func (o *OcctlRepository) Disconnect(c context.Context, username string) error {
	return o.oc.Disconnect(c, username)
}

func (o *OcctlRepository) ShowIPBans(c context.Context) (*[]occtl.IPBan, error) {
	ipBans, err := o.oc.ShowIPBans(c)
	if err != nil {
		return nil, err
	}
	return ipBans, nil
}

func (o *OcctlRepository) ShowIPBansPoint(c context.Context) (*[]occtl.IPBanPoints, error) {
	bansPoint, err := o.oc.ShowIPBansPoints(c)
	if err != nil {
		return nil, err
	}
	return bansPoint, nil
}

func (o *OcctlRepository) UnBanIP(c context.Context, ip string) error {
	return o.oc.UnBanIP(c, ip)
}

func (o *OcctlRepository) ShowStatus(c context.Context) string {
	status, err := o.oc.ShowStatus(c)
	if err != nil {
		return ""
	}
	return status
}
func (o *OcctlRepository) ShowIRoutes(c context.Context) (*[]occtl.IRoute, error) {
	iroutes, err := o.oc.ShowIRoutes(c)
	if err != nil {
		return nil, err
	}
	return iroutes, nil
}
func (o *OcctlRepository) ShowUser(c context.Context, username string) (*[]occtl.OcUser, error) {
	user, err := o.oc.ShowUser(c, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

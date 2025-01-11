package repository

import (
	"api/pkg/ocserv"
	"context"
)

type OcctlRepository struct {
	oc *ocserv.Handler
}

type OcctlRepositoryInterface interface {
	Reload(c context.Context) error
	OnlineUsers(c context.Context) (*[]ocserv.OcctlUser, error)
	Disconnect(c context.Context, username string) error
	ShowIPBans(c context.Context, points bool) *[]ocserv.IPBan
	UnBanIP(c context.Context, ip string) error
	ShowStatus(c context.Context) string
	ShowIRoutes(c context.Context) *[]ocserv.IRoute
	ShowUser(c context.Context, username string) (*[]ocserv.OcctlUser, error)
}

func NewOcctlRepository() *OcctlRepository {
	return &OcctlRepository{oc: ocserv.NewHandler()}
}

func (o *OcctlRepository) Reload(c context.Context) error {
	return o.oc.Occtl.Reload(c)
}

func (o *OcctlRepository) OnlineUsers(c context.Context) (*[]ocserv.OcctlUser, error) {
	return o.oc.Occtl.OnlineUsers(c)
}

func (o *OcctlRepository) Disconnect(c context.Context, username string) error {
	return o.oc.Occtl.Disconnect(c, username)
}

func (o *OcctlRepository) ShowIPBans(c context.Context, points bool) *[]ocserv.IPBan {
	return o.oc.Occtl.ShowIPBans(c, points)
}

func (o *OcctlRepository) UnBanIP(c context.Context, ip string) error {
	return o.oc.Occtl.UnBanIP(c, ip)
}

func (o *OcctlRepository) ShowStatus(c context.Context) string {
	return o.oc.Occtl.ShowStatus(c)
}
func (o *OcctlRepository) ShowIRoutes(c context.Context) *[]ocserv.IRoute {
	return o.oc.Occtl.ShowIRoutes(c)
}
func (o *OcctlRepository) ShowUser(c context.Context, username string) (*[]ocserv.OcctlUser, error) {
	return o.oc.Occtl.ShowUser(c, username)
}

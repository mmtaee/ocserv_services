package repository

import (
	"api/pkg/database"
	"api/pkg/ocserv"
	"context"
	"gorm.io/gorm"
)

type OcservGroupRepository struct {
	db *gorm.DB
	oc *ocserv.Handler
}

type OcservGroupRepositoryInterface interface {
	Groups(c context.Context) (*[]ocserv.OcGroupConfigInfo, error)
	GroupNames(c context.Context) (*[]string, error)
	UpdateDefaultGroup(context.Context, *ocserv.OcGroupConfig) error
	CreateOrUpdateGroup(context.Context, string, *ocserv.OcGroupConfig) error
	DeleteGroup(context.Context, string) error
}

func NewOcservGroupRepository() *OcservGroupRepository {
	return &OcservGroupRepository{
		db: database.Connection(),
		oc: ocserv.NewHandler(),
	}
}

func (o *OcservGroupRepository) Groups(c context.Context) (*[]ocserv.OcGroupConfigInfo, error) {
	groups, err := o.oc.Group.Groups(c)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (o *OcservGroupRepository) GroupNames(c context.Context) (*[]string, error) {
	groups, err := o.oc.Group.GroupNames(c)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (o *OcservGroupRepository) UpdateDefaultGroup(c context.Context, config *ocserv.OcGroupConfig) error {
	configMap := o.oc.ToMap(config)
	err := o.oc.Group.UpdateDefaultGroup(c, &configMap)
	if err != nil {
		return err
	}
	return o.oc.Occtl.Reload(c)
}

func (o *OcservGroupRepository) CreateOrUpdateGroup(c context.Context, name string, config *ocserv.OcGroupConfig) error {
	configMap := o.oc.ToMap(config)
	err := o.oc.Group.CreateOrUpdateGroup(c, name, &configMap)
	if err != nil {
		return err
	}
	return o.oc.Occtl.Reload(c)
}

func (o *OcservGroupRepository) DeleteGroup(c context.Context, name string) error {
	return o.oc.Group.DeleteGroup(c, name)
}

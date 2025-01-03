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
	UpdateDefaultGroup(context.Context, *ocserv.OcGroupConfig) error
	Groups(context.Context) (*[]ocserv.OcGroupConfig, error)
	CreateOrUpdateGroup(context.Context, string, *ocserv.OcGroupConfig) error
	DeleteGroup(context.Context, string) error
}

func NewOcservGroupRepository() *OcservGroupRepository {
	return &OcservGroupRepository{
		db: database.Connection(),
		oc: ocserv.NewHandler(),
	}
}

func (o *OcservGroupRepository) UpdateDefaultGroup(ctx context.Context, config *ocserv.OcGroupConfig) error {
	configMap := o.oc.ToMap(config)
	err := o.oc.Group.UpdateDefaultGroup(ctx, configMap)
	if err != nil {
		return err
	}
	return o.oc.Occtl.Reload()
}

func (o *OcservGroupRepository) Groups(ctx context.Context) (*[]ocserv.OcGroupConfig, error) {
	groups, err := o.oc.Group.Groups(ctx)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (o *OcservGroupRepository) CreateOrUpdateGroup(ctx context.Context, name string, config *ocserv.OcGroupConfig) error {
	configMap := o.oc.ToMap(config)
	err := o.oc.Group.CreateOrUpdateGroup(ctx, name, configMap)
	if err != nil {
		return err
	}
	return o.oc.Occtl.Reload()
}

func (o *OcservGroupRepository) DeleteGroup(ctx context.Context, name string) error {
	return o.oc.Group.DeleteGroup(ctx, name)
}

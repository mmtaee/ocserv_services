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
	UpdateDefaultGroup(context.Context, ocserv.OcGroupConfig) error
	Groups(context.Context) (*[]ocserv.OcGroupConfig, error)
	CreateGroup(context.Context) error
	UpdateGroup(context.Context) error
	DeleteGroup(context.Context) error
}

func NewOcservGroupRepository() *OcservGroupRepository {
	return &OcservGroupRepository{
		db: database.Connection(),
		oc: ocserv.NewHandler(),
	}
}

func (o *OcservGroupRepository) UpdateDefaultGroup(ctx context.Context, config ocserv.OcGroupConfig) error {
	configMap := o.oc.ToMap(config)
	err := o.oc.Group.UpdateDefaultGroup(ctx, configMap)
	if err != nil {
		return err
	}
	return o.oc.Occtl.Reload()
}

func (o *OcservGroupRepository) Groups(ctx context.Context) (*[]ocserv.OcGroupConfig, error) {
	return nil, nil
}

func (o *OcservGroupRepository) CreateGroup(ctx context.Context) error {
	//name := ctx.Value("name").(string)
	//config := ctx.Value("config").(models.GroupConfig)
	//configMap := o.oc.ToMap(config)
	return o.oc.Occtl.Reload()
}

func (o *OcservGroupRepository) UpdateGroup(ctx context.Context) error {
	//name := ctx.Value("name").(string)
	//config := ctx.Value("config").(models.GroupConfig)
	//configMap := o.oc.ToMap(config)
	return nil
}

func (o *OcservGroupRepository) DeleteGroup(ctx context.Context) error {
	//name := ctx.Value("name").(string)
	return nil
}

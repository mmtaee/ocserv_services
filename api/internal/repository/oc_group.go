package repository

import (
	"context"
	"encoding/json"
	"github.com/mmtaee/go-oc-utils/handler/occtl"
	"github.com/mmtaee/go-oc-utils/handler/ocgroup"
)

type OcservGroupRepository struct {
	ocGroup ocgroup.OcservGroupInterface
	occtl   occtl.OcInterface
}

type OcservGroupRepositoryInterface interface {
	Groups(c context.Context) (*[]ocgroup.OcservGroupConfigInfo, error)
	GroupNames(c context.Context) (*[]string, error)
	UpdateDefaultGroup(context.Context, *ocgroup.OcservGroupConfig) error
	CreateOrUpdateGroup(context.Context, string, *ocgroup.OcservGroupConfig) error
	DeleteGroup(context.Context, string) error
}

func NewOcservGroupRepository() *OcservGroupRepository {
	return &OcservGroupRepository{
		ocGroup: ocgroup.NewOcservGroup(),
		occtl:   occtl.NewOcctl(),
	}
}

func toMap(data interface{}) map[string]interface{} {
	b, _ := json.Marshal(&data)
	var dataStruct map[string]interface{}
	_ = json.Unmarshal(b, &dataStruct)
	return dataStruct
}

func (o *OcservGroupRepository) Groups(c context.Context) (*[]ocgroup.OcservGroupConfigInfo, error) {
	groups, err := o.ocGroup.List(c)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (o *OcservGroupRepository) GroupNames(c context.Context) (*[]string, error) {
	groups, err := o.ocGroup.NameList(c)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (o *OcservGroupRepository) UpdateDefaultGroup(c context.Context, config *ocgroup.OcservGroupConfig) error {
	configMap := toMap(config)
	err := o.ocGroup.UpdateDefault(c, &configMap)
	if err != nil {
		return err
	}
	return o.occtl.Reload(c)
}

func (o *OcservGroupRepository) CreateOrUpdateGroup(c context.Context, name string, config *ocgroup.OcservGroupConfig) error {
	configMap := toMap(config)
	err := o.ocGroup.Create(c, name, &configMap)
	if err != nil {
		return err
	}
	return o.occtl.Reload(c)
}

func (o *OcservGroupRepository) DeleteGroup(c context.Context, name string) error {
	if err := o.ocGroup.Delete(c, name); err != nil {
		return err
	}
	return o.occtl.Reload(c)
}

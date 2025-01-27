package repository

import (
	"api/pkg/event"
	"context"
	"encoding/json"
	"github.com/mmtaee/go-oc-utils/handler/occtl"
	"github.com/mmtaee/go-oc-utils/handler/ocgroup"
	"github.com/mmtaee/go-oc-utils/logger"
)

type OcservGroupRepository struct {
	ocGroup     ocgroup.OcservGroupInterface
	occtl       occtl.OcInterface
	WorkerEvent *event.WorkerEvent
}

type OcservGroupRepositoryInterface interface {
	Groups(c context.Context) (*[]ocgroup.OcservGroupConfigInfo, error)
	GroupNames(c context.Context) (*[]string, error)
	DefaultGroup(c context.Context) (*ocgroup.OcservGroupConfig, error)
	UpdateDefaultGroup(c context.Context, config *ocgroup.OcservGroupConfig) error
	CreateOrUpdateGroup(c context.Context, name string, config *ocgroup.OcservGroupConfig, create bool) error
	DeleteGroup(c context.Context, name string) error
}

func NewOcservGroupRepository() *OcservGroupRepository {
	return &OcservGroupRepository{
		ocGroup:     ocgroup.NewOcservGroup(),
		occtl:       occtl.NewOcctl(),
		WorkerEvent: event.GetWorker(),
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

func (o *OcservGroupRepository) DefaultGroup(c context.Context) (*ocgroup.OcservGroupConfig, error) {
	conf, err := o.ocGroup.DefaultGroup(c)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func (o *OcservGroupRepository) UpdateDefaultGroup(c context.Context, config *ocgroup.OcservGroupConfig) error {
	configMap := toMap(config)
	err := o.ocGroup.UpdateDefault(c, &configMap)
	if err != nil {
		return err
	}

	old, err := o.DefaultGroup(c)
	if err != nil {
		return nil
	}
	o.WorkerEvent.AddEvent(&event.SchemaEvent{
		EventType: "update_default_group",
		ModelName: "group",
		ModelUID:  "",
		OldState:  old,
		NewState:  config,
	})
	return o.occtl.Reload(c)
}

func (o *OcservGroupRepository) CreateOrUpdateGroup(c context.Context, name string, config *ocgroup.OcservGroupConfig, create bool) error {
	configMap := toMap(config)
	err := o.ocGroup.Create(c, name, &configMap)
	if err != nil {
		return err
	}

	var eventType string
	var oldState *ocgroup.OcservGroupConfig

	if create {
		eventType = "create_group"
		oldState = nil
	} else {
		eventType = "update_group"
		oldState, err = o.ocGroup.Group(c, name)
		if err != nil {
			logger.Logf(logger.ERROR, "get group err: %s", err.Error())
		}
	}

	o.WorkerEvent.AddEvent(&event.SchemaEvent{
		EventType: eventType,
		ModelName: "group",
		ModelUID:  name,
		OldState:  oldState,
		NewState:  config,
	})

	return o.occtl.Reload(c)
}

func (o *OcservGroupRepository) DeleteGroup(c context.Context, name string) error {
	if err := o.ocGroup.Delete(c, name); err != nil {
		return err
	}
	o.WorkerEvent.AddEvent(&event.SchemaEvent{
		EventType: "delete_group",
		ModelName: "group",
		ModelUID:  name,
		OldState:  nil,
		NewState:  nil,
	})
	return o.occtl.Reload(c)
}

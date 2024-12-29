package ocserv

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type OcGroup struct{}

// TODO: update all methods to support context for cancel

// OcGroupInterface all method in this interface need reload server config
type OcGroupInterface interface {
	UpdateDefaultGroup(map[string]interface{}) error
	CreateOrUpdateGroup(name string, config map[string]interface{}) error
	DeleteGroup(name string) error
}

func NewOcGroup() *OcGroup {
	return &OcGroup{}
}

func (o *OcGroup) UpdateDefaultGroup(config map[string]interface{}) error {
	file, err := os.Create(defaultGroup)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Printf("failed to close file: %v", closeErr)
		}
	}()
	for k, v := range config {
		if k == "dns" {
			dnsValues, ok := v.([]string)
			if !ok {
				return fmt.Errorf("invalid type for dns, expected []string but got %T", v)
			}
			for _, dns := range dnsValues {
				if _, err := file.WriteString(fmt.Sprintf("dns=%s\n", dns)); err != nil {
					return fmt.Errorf("failed to write to file: %w", err)
				}
			}
		} else {
			if _, err := file.WriteString(fmt.Sprintf("%s=%v\n", k, v)); err != nil {
				return fmt.Errorf("failed to write to file: %w", err)
			}
		}
	}
	return nil
}

func (o *OcGroup) CreateOrUpdateGroup(name string, config map[string]interface{}) error {
	file, err := os.Create(fmt.Sprintf("%s/%s", groupDir, name))
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Printf("failed to close file: %v", closeErr)
		}
	}()
	for k, v := range config {
		_, err = file.WriteString(fmt.Sprintf("%s=%s\n", k, v))
		if err != nil {
			break
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (o *OcGroup) DeleteGroup(name string) error {
	if name == "defaults" {
		return errors.New("default group cannot be deleted")
	}
	return os.Remove(fmt.Sprintf("%s/%s", groupDir, name))
}

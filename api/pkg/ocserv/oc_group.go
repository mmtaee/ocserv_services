package ocserv

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type OcGroup struct{}

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
	ch := make(chan error, 1)
	go func() {
		file, err := os.Create(defaultGroup)
		if err != nil {
			ch <- err
		}
		defer func(file *os.File) {
			err = file.Close()
			if err != nil {
				log.Println(err)
			}
		}(file)
		for k, v := range config {
			if k == "dns" {
				for _, dns := range v.([]string) {
					_, err = file.WriteString(fmt.Sprintf("dns=%s\n", dns))
				}
			} else {
				_, err = file.WriteString(fmt.Sprintf("%s=%s\n", k, v))
			}
		}
		ch <- err
	}()
	return <-ch
}

func (o *OcGroup) CreateOrUpdateGroup(name string, config map[string]interface{}) error {
	ch := make(chan error, 1)
	go func() {
		file, err := os.Create(fmt.Sprintf("%s/%s", groupDir, name))
		if err != nil {
			ch <- err
		}
		defer func(file *os.File) {
			err = file.Close()
			if err != nil {
				log.Println(err)
			}
		}(file)
		for k, v := range config {
			_, err = file.WriteString(fmt.Sprintf("%s=%s\n", k, v))
		}
		ch <- err
	}()
	return <-ch
}

func (o *OcGroup) DeleteGroup(name string) error {
	if name == "defaults" {
		return errors.New("default group cannot be deleted")
	}
	ch := make(chan error, 1)
	go func() {
		err := os.Remove(fmt.Sprintf("%s/%s", groupDir, name))
		ch <- err
	}()
	return <-ch
}

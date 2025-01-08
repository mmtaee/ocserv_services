package ocserv

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

type OcGroup struct{}

// OcGroupInterface all method in this interface need reload server config
type OcGroupInterface interface {
	Groups(c context.Context) (*[]OcGroupConfigInfo, error)
	GroupNames(c context.Context) (*[]string, error)
	UpdateDefaultGroup(c context.Context, config *map[string]interface{}) error
	CreateOrUpdateGroup(c context.Context, name string, config *map[string]interface{}) error
	DeleteGroup(context.Context, string) error
}

func NewOcGroup() *OcGroup {
	return &OcGroup{}
}

func (o *OcGroup) Groups(ctx context.Context) (*[]OcGroupConfigInfo, error) {
	var (
		result []OcGroupConfigInfo
		wg     sync.WaitGroup
	)
	err := WithContext(ctx, func() error {
		err := filepath.Walk(groupDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				result = append(result, OcGroupConfigInfo{
					Name: info.Name(),
					Path: path,
				})
			}
			return nil
		})
		if err != nil {
			return err
		}

		for i := range result {
			wg.Add(1)
			go func(data *OcGroupConfigInfo) {
				defer wg.Done()
				config, err := ParseConfFile(data.Path)
				if err != nil {
					fmt.Printf("Error parsing file %s: %v\n", data.Path, err)
					return
				}
				data.Config = config
			}(&result[i])
		}
		wg.Wait()
		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return &result, err
}

func (o *OcGroup) GroupNames(c context.Context) (*[]string, error) {
	var result []string
	err := WithContext(c, func() error {
		err := filepath.Walk(groupDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				result = append(result, info.Name())
			}
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(result)
	return &result, nil
}

func (o *OcGroup) UpdateDefaultGroup(ctx context.Context, config *map[string]interface{}) error {
	return WithContext(ctx, func() error {
		file, err := os.Create(defaultGroup)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)

		}
		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				log.Printf("failed to close file: %v", closeErr)
			}
		}()
		return GroupWriter(file, config)
	})
}

func (o *OcGroup) CreateOrUpdateGroup(ctx context.Context, name string, config *map[string]interface{}) error {
	return WithContext(ctx, func() error {
		file, err := os.Create(fmt.Sprintf("%s/%s", groupDir, name))
		if err != nil {
			return err
		}
		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				log.Printf("failed to close file: %v", closeErr)
			}
		}()
		return GroupWriter(file, config)
	})
}

func (o *OcGroup) DeleteGroup(ctx context.Context, name string) error {
	return WithContext(ctx, func() error {
		if name == "defaults" {
			return errors.New("default group cannot be deleted")
		}
		err := os.Remove(fmt.Sprintf("%s/%s", groupDir, name))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("group %s does not exist", name)
			}
			return fmt.Errorf("failed to delete group %s: %w", name, err)
		}
		return nil
	})
}

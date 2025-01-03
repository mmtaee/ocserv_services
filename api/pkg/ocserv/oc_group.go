package ocserv

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"sync"
)

type OcGroup struct{}

// OcGroupInterface all method in this interface need reload server config
type OcGroupInterface interface {
	Groups(ctx context.Context) (*[]OcGroupConfig, error)
	UpdateDefaultGroup(context.Context, map[string]interface{}) error
	CreateOrUpdateGroup(context.Context, string, map[string]interface{}) error
	DeleteGroup(context.Context, string) error
}

func NewOcGroup() *OcGroup {
	return &OcGroup{}
}

func withContext(ctx context.Context, operation func() error) error {
	done := make(chan error, 1)

	go func() {
		defer close(done)
		done <- operation()
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("operation canceled or timed out: %w", ctx.Err())
	case err := <-done:
		return err
	}
}

func (o *OcGroup) Groups(ctx context.Context) (*[]OcGroupConfig, error) {
	var (
		result     []OcGroupConfig
		groupFiles []string
		wg         sync.WaitGroup
	)
	ch := make(chan OcGroupConfig, len(groupFiles))
	errCh := make(chan error, len(groupFiles))

	err := withContext(ctx, func() error {
		err := filepath.Walk(groupDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				groupFiles = append(groupFiles, path)
			}
			return nil
		})
		if err != nil {
			return err
		}

		sort.Strings(groupFiles)

		for _, path := range groupFiles {
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				config, err := ParseConfFile(path)
				if err != nil {
					errCh <- fmt.Errorf("error parsing file %s: %v", path, err)
				}
				ch <- config
			}(path)
		}

		go func() {
			wg.Wait()
			close(ch)
			close(errCh)
		}()

		for {
			select {
			case config, ok := <-ch:
				if ok {
					result = append(result, config)
				}
			}
			if len(result) == len(groupFiles) && len(errCh) == 0 {
				break
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &result, err
}

func (o *OcGroup) UpdateDefaultGroup(ctx context.Context, config map[string]interface{}) error {
	return withContext(ctx, func() error {
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
	})
}

func (o *OcGroup) CreateOrUpdateGroup(ctx context.Context, name string, config map[string]interface{}) error {
	return withContext(ctx, func() error {
		file, err := os.Create(fmt.Sprintf("%s/%s", groupDir, name))
		if err != nil {
			return err
		}
		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				log.Printf("failed to close file: %v", closeErr)
			}
		}()

		val := reflect.ValueOf(config).Elem()
		typ := val.Type()
		for i := 0; i < val.NumField(); i++ {
			fieldName := typ.Field(i).Name
			fieldValue := val.Field(i)
			if !fieldValue.IsZero() {
				_, err = file.WriteString(fmt.Sprintf("%s=%v\n", fieldName, fieldValue.Interface()))
				if err != nil {
					break
				}
			}
		}

		if err != nil {
			return err
		}
		return nil
	})
}

func (o *OcGroup) DeleteGroup(ctx context.Context, name string) error {
	return withContext(ctx, func() error {
		if name == "defaults" {
			return errors.New("default group cannot be deleted")
		}
		return os.Remove(fmt.Sprintf("%s/%s", groupDir, name))
	})
}

package initialize

import (
	"api/pkg/config"
	"errors"
	"log"
	"os"
)

func checkSecret(secret string) error {
	if secret == "" {
		return errors.New("secret parameter is required")
	}
	file := config.GetApp().InitSecretFile
	_, err := os.Stat(file)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Println(err)
		}
		return err
	}
	content, err := os.ReadFile(config.GetApp().InitSecretFile)
	if err != nil {
		return nil
	}
	if secret != string(content) {
		return errors.New("invalid secret key or initial application preparation steps have already been completed")
	}
	return nil
}

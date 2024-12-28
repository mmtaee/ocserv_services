package password

import (
	"api/pkg/config"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func Create(password string) string {
	configs := config.GetApp()
	saltPassword := fmt.Sprintf("%s%s", password, configs.SecretKey)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(saltPassword), 12)
	if err != nil {
		return ""
	}
	return string(hashedPassword)
}

func Check(password, hashedPassword string) bool {
	configs := config.GetApp()
	saltPassword := fmt.Sprintf("%s%s", password, configs.SecretKey)
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(saltPassword))
	return err == nil
}

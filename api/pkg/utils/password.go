package utils

import (
	"api/pkg/config"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
)

type CustomPassword struct {
	Salt string
	Hash string
}

type CustomPasswordInterface interface {
	Salt(length int) string
	Create(passwd, salt string) string
	Check(passwd, hashedPassword, salt string) bool
}

func NewPassword(passwd string, saltLength ...int) *CustomPassword {
	length := 6
	if len(saltLength) > 0 {
		length = saltLength[0]
	}
	salt := createSalt(length)
	return &CustomPassword{
		Salt: salt,
		Hash: create(passwd, salt),
	}
}

func createSalt(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func create(passwd, salt string) string {
	secretKey := config.GetApp().SecretKey
	passwordHash := fmt.Sprintf("%s%s%s", salt, passwd, secretKey)
	hash := md5.New()
	hash.Write([]byte(passwordHash))
	return hex.EncodeToString(hash.Sum(nil))
}

func Check(passwd, hashedPassword, salt string) bool {
	hash := create(passwd, salt)
	if hashedPassword == hash {
		return true
	}
	return false
}

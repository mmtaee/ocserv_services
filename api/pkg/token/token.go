package token

import (
	"api/pkg/config"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"time"
)

func Create(userID uint, expire time.Time) string {
	cfg := config.GetApp()
	str := fmt.Sprintf("%d%v%s%v", time.Now().Unix(), userID, cfg.SecretKey, expire)
	hashed := sha256.New()
	hashed.Write([]byte(str))
	hash := hashed.Sum(nil)
	hashHex := hex.EncodeToString(hash)
	return hashHex
}

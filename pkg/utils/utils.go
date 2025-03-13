package utils

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

func GenerateSessionID() string {
	return hex.EncodeToString(md5.New().Sum([]byte(time.Now().String())))
}

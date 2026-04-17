package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// MD5 生成MD5哈希
func MD5(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

// SHA256 生成SHA256哈希
func SHA256(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}

// GenerateRandomHash 生成随机哈希
func GenerateRandomHash() string {
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())
	return MD5(timestamp)
}

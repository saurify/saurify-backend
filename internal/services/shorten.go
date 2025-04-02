package services

import (
	"crypto/sha256"
	"encoding/base64"
)

// func to create short hash for a known url
func GenerateShortCode(url string) string {
	hash := sha256.Sum256([]byte(url))
	return base64.URLEncoding.EncodeToString(hash[:])[:6]
}

package tools

import (
	"crypto/sha256"
	"encoding/hex"
)

// GenerateDeviceToken generates a unique token for a device
func GenerateDeviceToken(deviceID string) string {
	hash := sha256.New()
	hash.Write([]byte(deviceID))
	return hex.EncodeToString(hash.Sum(nil))
}

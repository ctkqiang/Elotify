package security

import (
	"pushnotification_services/internal/config"
)

var sharedEncryptionKey = []byte(config.JWECreds.KeyAES256)

func EncryptPayload(payload interface{}) (string, error) {
	if !config.JWECreds.Encrypt {
		return "", nil
	}

	// 如果加密未启用，返回空字符串
	return "", nil
}

func DecryptPayload(jweString string, dest interface{}) error {
	if !config.JWECreds.Encrypt {
		return nil
	}

	return nil
}

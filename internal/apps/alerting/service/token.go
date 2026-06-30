package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func generateCallbackToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate callback token: %w", err)
	}
	return hex.EncodeToString(buf), nil
}

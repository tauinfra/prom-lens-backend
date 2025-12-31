package encryption

import (
	"encoding/base64"
	"fmt"
	"valyria-backend/internal/core/configs"

	"crypto/rand"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
)

type Encryptor interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(encrypted string) (string, error)
}

type SecretBoxEncryptor struct {
	key *[32]byte
}

func NewSecretBoxEncryptor(cfg *configs.Config) (Encryptor, error) {
	key := &[32]byte{}
	copy(key[:], []byte(cfg.App.EncryptionKey))
	return &SecretBoxEncryptor{key: key}, nil
}

func (e *SecretBoxEncryptor) Encrypt(plaintext string) (string, error) {
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return "", fmt.Errorf("generate nonce failed: %w", err)
	}

	encrypted := secretbox.Seal(nonce[:], []byte(plaintext), &nonce, e.key)
	return base64.URLEncoding.EncodeToString(encrypted), nil
}

func (e *SecretBoxEncryptor) Decrypt(encrypted string) (string, error) {
	// 解码 base64
	data, err := base64.URLEncoding.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("decode base64 failed: %w", err)
	}
	// 检查数据长度
	if len(data) < 24 {
		return "", fmt.Errorf("encrypted data too short")
	}
	var nonce [24]byte
	copy(nonce[:], data[:24])

	decrypted, ok := secretbox.Open(nil, data[24:], &nonce, e.key)
	if !ok {
		return "", fmt.Errorf("decryption failed")
	}

	return string(decrypted), nil
}

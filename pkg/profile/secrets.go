package profile

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

type Secrets interface {
	Encrypt(ctx context.Context, plain string) ([]byte, error)
	Decrypt(ctx context.Context, encrypted []byte) (string, error)
}

func NewSecrets(config Config) Secrets {
	return &secrets{config: config}
}

type secrets struct {
	config Config
}

func (s *secrets) Encrypt(ctx context.Context, plaintext string) ([]byte, error) {
	key := []byte(s.config.Passphrase)
	if len(key) != 32 {
		return nil, errors.New("key is not 32 bytes")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, []byte(plaintext), nil), nil
}

func (s *secrets) Decrypt(ctx context.Context, ciphertext []byte) (string, error) {
	block, err := aes.NewCipher([]byte(s.config.Passphrase))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return "", errors.New("malformed ciphertext")
	}

	plaintext, err := gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

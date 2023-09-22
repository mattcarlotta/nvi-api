package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

var ENCRYPTION_KEY = []byte(GetEnv("ENCRYPTION_KEY"))

func CreateEncryptedSecretValue(plaintext []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(ENCRYPTION_KEY)
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	return ciphertext, nonce, nil
}

func DecryptSecretValue(data []byte, nonce []byte) ([]byte, error) {
	if len(nonce) == 0 {
		return nil, errors.New("the provided nonce is not valid because it has no length")
	}

	block, err := aes.NewCipher(ENCRYPTION_KEY)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return gcm.Open(nil, nonce, data, nil)
}

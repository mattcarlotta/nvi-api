package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

var PASSPHRASE = []byte(GetEnv("PASSPHRASE"))

func CreateEncryptedSecretValue(plaintext []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(PASSPHRASE)
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
	block, err := aes.NewCipher(PASSPHRASE)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return gcm.Open(nil, nonce, data, nil)
}

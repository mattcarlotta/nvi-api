package utils

import "golang.org/x/crypto/bcrypt"

func ComparePassword(currentPassword []byte, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(currentPassword, password)
	return err == nil
}

func CreateEncryptedPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, 0)
}

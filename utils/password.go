package utils

import "golang.org/x/crypto/bcrypt"

func CompareEncryptedText(encyrptedText []byte, text []byte) bool {
	err := bcrypt.CompareHashAndPassword(encyrptedText, text)
	return err == nil
}

func CreateEncryptedText(text []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(text, 0)
}

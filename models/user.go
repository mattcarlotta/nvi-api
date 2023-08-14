package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Email     string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password  []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (user *User) setPassword(password string) {
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// TODO(carlotta): handle potential error
	}
	user.Password = pwd
}

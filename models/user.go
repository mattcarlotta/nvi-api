package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Email     string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password  []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

// type NewUser struct {
// 	Name     string `json:"name"`
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

// func (user *User) SetPassword(password string) {
// 	pwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		// TODO(carlotta): handle potential error
// 	}
// 	user.Password = pwd
// }

func (user *User) MatchPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	return err == nil
}

func (user *User) BeforeSave(tx *gorm.DB) (err error) {
	if pw, err := bcrypt.GenerateFromPassword(user.Password, 0); err == nil {
		tx.Statement.SetColumn("Password", pw)
	}
	return
}

package models

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/mattcarlotta/nvi-api/utils"
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

func (user *User) MatchPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	return err == nil
}

func (user *User) GenerateSessionToken() (time.Time, string, error) {
	exp := time.Now().Add(time.Hour * 24 * 30)
	claims := &utils.JWTSessionClaim{
		Email:  user.Email,
		Name:   user.Name,
		UserId: user.ID.String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(utils.JWT_SECRET_KEY)
	return exp, tokenString, err
}

func (user *User) BeforeSave(tx *gorm.DB) (err error) {
	if pw, err := bcrypt.GenerateFromPassword(user.Password, 0); err == nil {
		tx.Statement.SetColumn("Password", pw)
	}
	return nil
}

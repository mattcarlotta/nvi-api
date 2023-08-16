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
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Email     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password  []byte    `gorm:"not null" json:"password"`
	Verified  bool      `gorm:"default:false" json:"verified"`
	Token     *[]byte   `gorm:"default:null" json:"token"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (user *User) MatchPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	return err == nil
}

func (user *User) GenerateSessionToken() (string, time.Time, error) {
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
	return tokenString, exp, err
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	if pw, err := bcrypt.GenerateFromPassword(user.Password, 0); err == nil {
		tx.Statement.SetColumn("Password", pw)
	}
	return nil
}

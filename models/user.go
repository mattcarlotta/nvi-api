package models

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/mattcarlotta/nvi-api/utils"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Email     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password  []byte    `gorm:"not null" json:"-"`
	Verified  bool      `gorm:"default:false" json:"-"`
	Token     *[]byte   `gorm:"default:null" json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (user *User) MatchPassword(password string) bool {
	return utils.CompareEncryptedText(user.Password, []byte(password))
}

func (user *User) GenerateSessionToken() (string, time.Time, error) {
	exp := time.Now().Add(time.Hour * 24 * 30)
	claims := &utils.JWTSessionClaim{
		Email:  user.Email,
		Name:   user.Name,
		UserID: user.ID.String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(utils.JWT_SECRET_KEY)
	return tokenString, exp, err
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	if pw, err := utils.CreateEncryptedText(user.Password); err == nil {
		tx.Statement.SetColumn("Password", pw)
	}
	return nil
}

type ReqRegisterUser struct {
	Name     string `json:"name" validate:"required,gte=2,lte=255"`
	Email    string `json:"email" validate:"required,email,lte=100"`
	Password string `json:"password" validate:"required,gte=5,lte=36"`
}

type ReqLoginUser struct {
	Email    string `json:"email" validate:"required,email,lte=100"`
	Password string `json:"password" validate:"required,gte=5,lte=36"`
}

type ReqUpdateUser struct {
	Password string `json:"password" validate:"required,gte=5,lte=36"`
	Token    string `json:"token" validate:"required"`
}

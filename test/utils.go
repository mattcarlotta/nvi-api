package testutils

import (
	"log"

	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
)

var Password = []byte("password123")
var NewToken = []byte("hello")
var StrPassword = string(Password)

type TestResponse struct {
	Route        string
	Method       string
	ExpectedCode int
}

func DeleteUser(email *string) {
	db := database.GetConnection()
	var existingUser models.User
	if err := db.Where(&models.User{Email: *email}).First(&existingUser).Error; err != nil {
		log.Fatal("unable to locate created user")
	}

	if err := db.Delete(&existingUser).Error; err != nil {
		log.Fatal("unable to delete created user")
	}

}

func CreateUser(email *string, verified bool) {
	db := database.GetConnection()

	newUser := &models.User{
		Name:     "Name",
		Email:    *email,
		Password: Password,
		Token:    &NewToken,
		Verified: verified,
	}

	if err := db.Create(newUser).Error; err != nil {
		log.Fatal("unable to create a user")
	}
}

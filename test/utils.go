package testutils

import (
	"encoding/json"
	"io"
	"log"

	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
	"github.com/mattcarlotta/nvi-api/utils"
)

var StrPassword = "password123"
var Password = []byte(StrPassword)

type TestResponse struct {
	Route        string
	Method       string
	ExpectedCode int
}

func DeleteUser(existingUser *models.User) {
	db := database.GetConnection()
	if err := db.Delete(&existingUser).Error; err != nil {
		log.Fatal("unable to delete created user")
	}
}

func RemoveUserByEmail(email string) {
	db := database.GetConnection()
	var existingUser models.User
	if err := db.Where(&models.User{Email: email}).First(&existingUser).Error; err != nil {
		log.Fatal("unable to locate created user")
	}

	if err := db.Delete(&existingUser).Error; err != nil {
		log.Fatal("unable to delete created user")
	}

}

func CreateUser(email string, verified bool) (models.User, string) {
	db := database.GetConnection()

	token, _, err := utils.GenerateUserToken(email)
	if err != nil {
		log.Fatal("unable to generate a new user token")
	}

	tokenByte := []byte(token)
	newUser := &models.User{
		Name:     "Name",
		Email:    email,
		Password: Password,
		Token:    &tokenByte,
		Verified: verified,
	}

	if err := db.Create(newUser).Error; err != nil {
		log.Fatal("unable to create a user")
	}

	var existingUser models.User
	if err := db.Where(&models.User{Email: email}).First(&existingUser).Error; err != nil {
		log.Fatal("unable to locate created user")
	}

	token, _, err = existingUser.GenerateSessionToken()
	if err != nil {
		log.Fatal("Unable to generate a user session token")
	}

	return existingUser, token
}

// func DeleteEnvironment(existingEnvironment *models.Environment) {
// 	db := database.GetConnection()
// 	if err := db.Delete(&existingEnvironment).Error; err != nil {
// 		log.Fatal("unable to delete created environment")
// 	}
// }

func CreateEnvironment(envName string, userSessionID string) models.Environment {
	db := database.GetConnection()

	token, err := utils.ValidateSessionToken(userSessionID)
	if err != nil {
		log.Fatal("unable to parse session token")
	}

	parsedID, err := utils.ParseUUID(token.UserID)
	if err != nil {
		log.Fatal("unable to parse user token id from session")
	}

	newEnv := models.Environment{Name: envName, UserID: parsedID}
	if err := db.Create(&newEnv).Error; err != nil {
		log.Fatal("unable to create the new environment")
	}

	var environment models.Environment
	if err := db.Where(
		&models.Environment{Name: newEnv.Name, UserID: parsedID},
	).First(&environment).Error; err != nil {
		log.Fatal("unable to locate the new environment")
	}

	return environment
}

// func ParseJSONSuccessBody(body *io.ReadCloser) utils.ResponseError {
// 	var res utils.ResponseError
// 	responseBodyBytes, _ := io.ReadAll(*body)
// 	_ = json.Unmarshal(responseBodyBytes, &res)

// 	return res
// }

func ParseJSONErrorBody(body *io.ReadCloser) utils.ResponseError {
	var errResponse utils.ResponseError
	responseBodyBytes, _ := io.ReadAll(*body)
	_ = json.Unmarshal(responseBodyBytes, &errResponse)

	return errResponse
}

func ParseText(body *io.ReadCloser) string {
	resBody, _ := io.ReadAll(*body)
	return string(resBody)
}
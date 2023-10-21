package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
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

func CreateUser(email string, verified bool) (models.User, string, string) {
	db := database.GetConnection()

	authToken, _, err := utils.GenerateUserToken(email)
	if err != nil {
		log.Fatal("unable to generate a new user token")
	}

	newUser := models.User{
		Name:     "Name",
		Email:    email,
		Password: Password,
		Verified: verified,
	}

	if err := db.Create(&newUser).Error; err != nil {
		log.Fatalf("unable to create a user: %v", err)
	}

	token, _, err := newUser.GenerateSessionToken()
	if err != nil {
		log.Fatalf("unable to generate a user session token: %v", err)
	}

	return newUser, token, authToken
}

// func DeleteEnvironment(existingEnvironment *models.Environment) {
// 	db := database.GetConnection()
// 	if err := db.Delete(&existingEnvironment).Error; err != nil {
// 		log.Fatal("unable to delete created environment")
// 	}
// }

func ParseSessionId(userSessionID string) uuid.UUID {
	token, err := utils.ValidateSessionToken(userSessionID)
	if err != nil {
		log.Fatalf("unable to parse session token: %v", err)
	}

	parsedID, err := utils.ParseUUID(token.UserID)
	if err != nil {
		log.Fatal("unable to parse user token id from session")
	}

	return parsedID
}

func CreateProject(name string, userSessionID string) models.Project {
	db := database.GetConnection()

	parsedID := ParseSessionId(userSessionID)

	newProject := models.Project{Name: name, UserID: parsedID}
	if err := db.Create(&newProject).Error; err != nil {
		log.Fatalf("unable to create a new project %s: %v", name, err)
	}

	return newProject

}

func CreateEnvironment(envName string, projectID uuid.UUID, userSessionID string) models.Environment {
	db := database.GetConnection()

	parsedUserSessionID := ParseSessionId(userSessionID)

	newEnv := models.Environment{Name: envName, ProjectID: projectID, UserID: parsedUserSessionID}
	if err := db.Create(&newEnv).Error; err != nil {
		log.Fatalf("unable to create the new environment: %v", err)
	}

	return newEnv
}

func CreateEnvironmentAndSecret(envName string, projectID uuid.UUID, secretKey string, secretValue string, userSessionID string) (models.Environment, models.Secret) {
	newEnv := CreateEnvironment(envName, projectID, userSessionID)
	db := database.GetConnection()

	parsedID := ParseSessionId(userSessionID)

	environments := []models.Environment{newEnv}
	secret := models.Secret{Key: secretKey, Value: []byte(secretValue), UserID: parsedID, Environments: environments}
	if err := db.Create(&secret).Error; err != nil {
		log.Fatalf("unable to create new environment: %v", err)
	}

	return newEnv, secret
}

func CreateProjectAndEnvironmentAndSecret(projectName string, envName string, secretKey string, secretValue string, userSessionID string) (models.Project, models.Environment, models.Secret) {
	newProject := CreateProject(projectName, userSessionID)
	newEnv, newSecret := CreateEnvironmentAndSecret(envName, newProject.ID, secretKey, secretValue, userSessionID)

	return newProject, newEnv, newSecret
}

func CreateHTTPRequest(test *TestResponse, body ...interface{}) *http.Request {
	var bodyBuf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&bodyBuf).Encode(body[0]); err != nil {
			log.Fatal(err)
		}
	} else {
		bodyBuf.Write(nil)
	}

	req := httptest.NewRequest(test.Method, test.Route, &bodyBuf)
	req.Header.Add("Content-Type", "application/json")

	return req
}

func CreateAuthHTTPRequest(test *TestResponse, token *string, body ...interface{}) *http.Request {
	req := CreateHTTPRequest(test, body...)
	req.Header.Add("Cookie", fmt.Sprintf("SESSION_TOKEN=%s", *token))

	return req
}

// func ParseJSONSuccessBody(body *io.ReadCloser) utils.ResponseError {
// 	var res utils.ResponseError
// 	responseBodyBytes, _ := io.ReadAll(*body)
// 	_ = json.Unmarshal(responseBodyBytes, &res)

//		return res
//	}

func ParseJSONBodyError(body *io.ReadCloser) utils.ResponseError {
	var errResponse utils.ResponseError
	if err := json.NewDecoder(*body).Decode(&errResponse); err != nil {
		log.Fatal(err)
	}

	return errResponse
}

func ParseText(body *io.ReadCloser) string {
	resBody, _ := io.ReadAll(*body)
	return string(resBody)
}

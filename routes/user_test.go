package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
	"github.com/mattcarlotta/nvi-api/utils"
	"github.com/stretchr/testify/assert"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func TestRegisterUserEmptyBody(t *testing.T) {
	test := struct {
		route        string
		method       string
		expectedCode int
	}{
		route:        "/register",
		method:       fiber.MethodPost,
		expectedCode: fiber.StatusBadRequest,
	}

	req := httptest.NewRequest(test.method, test.route, nil)
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to register user controller")
	}

	var errResponse ErrorResponse
	responseBodyBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(responseBodyBytes, &errResponse)

	assert.Equal(t, test.expectedCode, resp.StatusCode)
	assert.Equal(t, errResponse.Error, utils.ErrorCode[utils.RegisterEmptyBody])
}

func TestRegisterUserInvalidBody(t *testing.T) {
	user := &models.ReqRegisterUser{
		Name:     "invalid",
		Email:    "invalidexample", // invalid email to trigger validation failure
		Password: []byte("password123"),
	}

	test := struct {
		route        string
		method       string
		expectedCode int
	}{
		route:        "/register",
		method:       fiber.MethodPost,
		expectedCode: fiber.StatusBadRequest,
	}

	reqBodyStr, _ := json.Marshal(user)
	req := httptest.NewRequest(test.method, test.route, bytes.NewBufferString(string(reqBodyStr)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to register user controller")
	}

	var errResponse ErrorResponse
	responseBodyBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(responseBodyBytes, &errResponse)

	assert.Equal(t, test.expectedCode, resp.StatusCode)
	assert.Equal(t, errResponse.Error, utils.ErrorCode[utils.RegisterInvalidBody])
}

func TestRegisterEmailTaken(t *testing.T) {
	db := database.GetConnection()

	password := []byte("password123")
	newToken := []byte("hello")
	email := "taken_email@example.com"
	newUser := &models.User{
		Name:     "Taken",
		Email:    email,
		Password: password,
		Token:    &newToken,
	}

	if err := db.Create(newUser).Error; err != nil {
		log.Fatal("unable to create a user")
	}

	user := &models.ReqRegisterUser{
		Name:     "Taken",
		Email:    email,
		Password: []byte("password123"),
	}

	test := struct {
		route        string
		method       string
		expectedCode int
	}{
		route:        "/register",
		method:       fiber.MethodPost,
		expectedCode: fiber.StatusOK,
	}

	reqBodyStr, _ := json.Marshal(user)
	req := httptest.NewRequest(test.method, test.route, bytes.NewBufferString(string(reqBodyStr)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to register user controller")
	}

	var errResponse ErrorResponse
	responseBodyBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(responseBodyBytes, &errResponse)

	defer func() {
		var existingUser models.User
		if err := db.Where(&models.User{Email: newUser.Email}).First(&existingUser).Error; err != nil {
			log.Fatal("unable to locate registered user")
		}

		if err := db.Delete(&existingUser).Error; err != nil {
			log.Fatal("unable to delete created user")
		}
	}()

	assert.Equal(t, test.expectedCode, resp.StatusCode)
	assert.Equal(t, errResponse.Error, utils.ErrorCode[utils.RegisterEmailTaken])
}

func TestRegisterUserSuccess(t *testing.T) {
	db := database.GetConnection()

	user := &models.ReqRegisterUser{
		Name:     "Register",
		Email:    "registeruser@example.com",
		Password: []byte("password123"),
	}

	test := struct {
		route        string
		method       string
		expectedCode int
	}{
		route:        "/register",
		method:       fiber.MethodPost,
		expectedCode: fiber.StatusCreated,
	}

	reqBodyStr, _ := json.Marshal(user)
	req := httptest.NewRequest(test.method, test.route, bytes.NewBufferString(string(reqBodyStr)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to register user controller")
	}

	body, _ := io.ReadAll(resp.Body)
	bodyMessage := string(body)

	defer func() {
		var existingUser models.User
		if err := db.Where(&models.User{Email: user.Email}).First(&existingUser).Error; err != nil {
			log.Fatal("unable to locate registered user")
		}

		if err := db.Delete(&existingUser).Error; err != nil {
			log.Fatal("unable to delete registered user")
		}
	}()

	assert.Equal(t, test.expectedCode, resp.StatusCode)
	assert.Equal(t, bodyMessage, fmt.Sprintf(
		"Welcome, %s! Please check your %s inbox for steps to verify your account.", user.Name, user.Email,
	))
}

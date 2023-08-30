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

func cleanup(email string) {
	db := database.GetConnection()

	var existingUser models.User
	if err := db.Where(&models.User{Email: email}).First(&existingUser).Error; err != nil {
		log.Fatal("unable to locate created user")
	}

	if err := db.Delete(&existingUser).Error; err != nil {
		log.Fatal("unable to delete created user")
	}
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

	var errResponse utils.ResponseError
	responseBodyBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(responseBodyBytes, &errResponse)

	assert.Equal(t, test.expectedCode, resp.StatusCode)
	assert.Equal(t, errResponse.Error, utils.ErrorCode[utils.RegisterEmptyBody])
}

func TestRegisterUserInvalidBody(t *testing.T) {
	user := &models.ReqRegisterUser{
		Name:     "invalid",
		Email:    "invalidexample", // invalid email to trigger validation failure
		Password: "password123",
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

	var errResponse utils.ResponseError
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
		Password: string(password),
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

	var errResponse utils.ResponseError
	responseBodyBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(responseBodyBytes, &errResponse)

	defer cleanup(newUser.Email)

	assert.Equal(t, test.expectedCode, resp.StatusCode)
	assert.Equal(t, errResponse.Error, utils.ErrorCode[utils.RegisterEmailTaken])
}

func TestRegisterUserSuccess(t *testing.T) {
	user := &models.ReqRegisterUser{
		Name:     "Register",
		Email:    "registeruser@example.com",
		Password: "password123",
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

	defer cleanup(user.Email)

	assert.Equal(t, test.expectedCode, resp.StatusCode)
	assert.Equal(t, bodyMessage, fmt.Sprintf(
		"Welcome, %s! Please check your %s inbox for steps to verify your account.", user.Name, user.Email,
	))
}

func TestLoginUserEmptyBody(t *testing.T) {
	test := struct {
		route        string
		method       string
		expectedCode int
	}{
		route:        "/login",
		method:       fiber.MethodPost,
		expectedCode: fiber.StatusBadRequest,
	}

	req := httptest.NewRequest(test.method, test.route, nil)
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to user login controller")
	}

	var errResponse utils.ResponseError
	responseBodyBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(responseBodyBytes, &errResponse)

	assert.Equal(t, test.expectedCode, resp.StatusCode)
	assert.Equal(t, errResponse.Error, utils.ErrorCode[utils.LoginEmptyBody])
}

func TestLoginUserInvalidBody(t *testing.T) {
	user := &models.ReqLoginUser{
		Email:    "invalidexample", // invalid email to trigger validation failure
		Password: "password123",
	}

	test := struct {
		route        string
		method       string
		expectedCode int
	}{
		route:        "/login",
		method:       fiber.MethodPost,
		expectedCode: fiber.StatusBadRequest,
	}

	reqBodyStr, _ := json.Marshal(user)
	req := httptest.NewRequest(test.method, test.route, bytes.NewBufferString(string(reqBodyStr)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to login user controller")
	}

	var errResponse utils.ResponseError
	responseBodyBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(responseBodyBytes, &errResponse)

	assert.Equal(t, test.expectedCode, resp.StatusCode)
	assert.Equal(t, errResponse.Error, utils.ErrorCode[utils.LoginInvalidBody])
}

func TestLoginUnregisteredEmail(t *testing.T) {
	user := &models.ReqLoginUser{
		Email:    "non_existent_email@example.com",
		Password: "password123",
	}

	test := struct {
		route        string
		method       string
		expectedCode int
	}{
		route:        "/login",
		method:       fiber.MethodPost,
		expectedCode: fiber.StatusOK,
	}

	reqBodyStr, _ := json.Marshal(user)
	req := httptest.NewRequest(test.method, test.route, bytes.NewBufferString(string(reqBodyStr)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to login user controller")
	}

	var errResponse utils.ResponseError
	responseBodyBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(responseBodyBytes, &errResponse)

	assert.Equal(t, test.expectedCode, resp.StatusCode)
	assert.Equal(t, errResponse.Error, utils.ErrorCode[utils.LoginUnregisteredEmail])
}

func TestLoginInvalidPassword(t *testing.T) {
	db := database.GetConnection()

	password := []byte("password123")
	newToken := []byte("hello")
	email := "login_invalid_password@example.com"
	newUser := &models.User{
		Name:     "Login",
		Email:    email,
		Password: password,
		Token:    &newToken,
	}

	if err := db.Create(newUser).Error; err != nil {
		log.Fatal("unable to create a user")
	}

	badPassword := append(password, []byte("4")...)
	user := &models.ReqLoginUser{
		Email:    email,
		Password: string(badPassword),
	}

	test := struct {
		route        string
		method       string
		expectedCode int
	}{
		route:        "/login",
		method:       fiber.MethodPost,
		expectedCode: fiber.StatusOK,
	}

	reqBodyStr, _ := json.Marshal(user)
	req := httptest.NewRequest(test.method, test.route, bytes.NewBufferString(string(reqBodyStr)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to login user controller")
	}

	var errResponse utils.ResponseError
	responseBodyBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(responseBodyBytes, &errResponse)

	defer cleanup(newUser.Email)

	assert.Equal(t, test.expectedCode, resp.StatusCode)
	assert.Equal(t, errResponse.Error, utils.ErrorCode[utils.LoginInvalidPassword])
}

func TestLoginAccountNotVerified(t *testing.T) {
	db := database.GetConnection()

	password := []byte("password123")
	newToken := []byte("hello")
	email := "login_account_not_verified@example.com"
	newUser := &models.User{
		Name:     "Login",
		Email:    email,
		Password: password,
		Token:    &newToken,
	}

	if err := db.Create(newUser).Error; err != nil {
		log.Fatal("unable to create a user")
	}

	user := &models.ReqLoginUser{
		Email:    email,
		Password: string(password),
	}

	test := struct {
		route        string
		method       string
		expectedCode int
	}{
		route:        "/login",
		method:       fiber.MethodPost,
		expectedCode: fiber.StatusUnauthorized,
	}

	reqBodyStr, _ := json.Marshal(user)
	req := httptest.NewRequest(test.method, test.route, bytes.NewBufferString(string(reqBodyStr)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to login user controller")
	}

	var errResponse utils.ResponseError
	responseBodyBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(responseBodyBytes, &errResponse)

	defer cleanup(newUser.Email)

	assert.Equal(t, test.expectedCode, resp.StatusCode)
	assert.Equal(t, errResponse.Error, utils.ErrorCode[utils.LoginAccountNotVerified])
}

func TestLoginSuccess(t *testing.T) {
	db := database.GetConnection()

	password := []byte("password123")
	newToken := []byte("hello")
	email := "login_success@example.com"
	newUser := &models.User{
		Name:     "Login",
		Email:    email,
		Password: password,
		Token:    &newToken,
		Verified: true,
	}

	if err := db.Create(newUser).Error; err != nil {
		log.Fatal("unable to create a user")
	}

	user := &models.ReqLoginUser{
		Email:    email,
		Password: string(password),
	}

	test := struct {
		route        string
		method       string
		expectedCode int
	}{
		route:        "/login",
		method:       fiber.MethodPost,
		expectedCode: fiber.StatusOK,
	}

	reqBodyStr, _ := json.Marshal(user)
	req := httptest.NewRequest(test.method, test.route, bytes.NewBufferString(string(reqBodyStr)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to login user controller")
	}

	cookie := resp.Header.Get("Set-Cookie")

	defer cleanup(newUser.Email)

	assert.Equal(t, test.expectedCode, resp.StatusCode)
	assert.NotEmpty(t, cookie)
}

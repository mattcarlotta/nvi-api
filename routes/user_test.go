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
	"github.com/mattcarlotta/nvi-api/models"
	"github.com/mattcarlotta/nvi-api/test"
	"github.com/mattcarlotta/nvi-api/utils"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUserEmptyBody(t *testing.T) {
	test := &testutils.TestResponse{
		Route:        "/register",
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusBadRequest,
	}

	req := httptest.NewRequest(test.Method, test.Route, nil)
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to register user controller")
	}

	var errResponse utils.ResponseError
	responseBodyBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(responseBodyBytes, &errResponse)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, errResponse.Error, utils.ErrorCode[utils.RegisterEmptyBody])
}

func TestRegisterUserInvalidBody(t *testing.T) {
	user := &models.ReqRegisterUser{
		Name:     "invalid",
		Email:    "invalidexample", // invalid email to trigger validation failure
		Password: string(testutils.Password),
	}

	test := &testutils.TestResponse{
		Route:        "/register",
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusBadRequest,
	}

	reqBodyStr, _ := json.Marshal(user)
	req := httptest.NewRequest(test.Method, test.Route, bytes.NewBufferString(string(reqBodyStr)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to register user controller")
	}

	var errResponse utils.ResponseError
	responseBodyBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(responseBodyBytes, &errResponse)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, errResponse.Error, utils.ErrorCode[utils.RegisterInvalidBody])
}

func TestRegisterEmailTaken(t *testing.T) {
	email := "taken_email@example.com"
	testutils.CreateUser(&email, false)

	user := &models.ReqRegisterUser{
		Name:     "Taken",
		Email:    email,
		Password: testutils.StrPassword,
	}

	test := &testutils.TestResponse{
		Route:        "/register",
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusOK,
	}

	reqBodyStr, _ := json.Marshal(user)
	req := httptest.NewRequest(test.Method, test.Route, bytes.NewBufferString(string(reqBodyStr)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to register user controller")
	}

	var errResponse utils.ResponseError
	responseBodyBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(responseBodyBytes, &errResponse)

	defer testutils.DeleteUser(&email)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, errResponse.Error, utils.ErrorCode[utils.RegisterEmailTaken])
}

func TestRegisterUserSuccess(t *testing.T) {
	user := &models.ReqRegisterUser{
		Name:     "Register",
		Email:    "registeruser@example.com",
		Password: testutils.StrPassword,
	}

	test := &testutils.TestResponse{
		Route:        "/register",
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusCreated,
	}

	reqBodyStr, _ := json.Marshal(user)
	req := httptest.NewRequest(test.Method, test.Route, bytes.NewBufferString(string(reqBodyStr)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to register user controller")
	}

	body, _ := io.ReadAll(resp.Body)
	bodyMessage := string(body)

	defer testutils.DeleteUser(&user.Email)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, bodyMessage, fmt.Sprintf(
		"Welcome, %s! Please check your %s inbox for steps to verify your account.", user.Name, user.Email,
	))
}

func TestLoginUserEmptyBody(t *testing.T) {
	test := &testutils.TestResponse{
		Route:        "/login",
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusBadRequest,
	}

	req := httptest.NewRequest(test.Method, test.Route, nil)
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to user login controller")
	}

	var errResponse utils.ResponseError
	responseBodyBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(responseBodyBytes, &errResponse)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, errResponse.Error, utils.ErrorCode[utils.LoginEmptyBody])
}

func TestLoginUserInvalidBody(t *testing.T) {
	user := &models.ReqLoginUser{
		Email:    "invalidexample", // invalid email to trigger validation failure
		Password: testutils.StrPassword,
	}

	test := &testutils.TestResponse{
		Route:        "/login",
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusBadRequest,
	}

	reqBodyStr, _ := json.Marshal(user)
	req := httptest.NewRequest(test.Method, test.Route, bytes.NewBufferString(string(reqBodyStr)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to login user controller")
	}

	var errResponse utils.ResponseError
	responseBodyBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(responseBodyBytes, &errResponse)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, errResponse.Error, utils.ErrorCode[utils.LoginInvalidBody])
}

func TestLoginUnregisteredEmail(t *testing.T) {
	user := &models.ReqLoginUser{
		Email:    "non_existent_email@example.com",
		Password: testutils.StrPassword,
	}

	test := &testutils.TestResponse{
		Route:        "/login",
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusOK,
	}

	reqBodyStr, _ := json.Marshal(user)
	req := httptest.NewRequest(test.Method, test.Route, bytes.NewBufferString(string(reqBodyStr)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to login user controller")
	}

	var errResponse utils.ResponseError
	responseBodyBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(responseBodyBytes, &errResponse)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, errResponse.Error, utils.ErrorCode[utils.LoginUnregisteredEmail])
}

func TestLoginInvalidPassword(t *testing.T) {
	email := "login_invalid_password@example.com"
	testutils.CreateUser(&email, false)

	badPassword := append(testutils.Password, []byte("4")...)
	user := &models.ReqLoginUser{
		Email:    email,
		Password: string(badPassword),
	}

	test := &testutils.TestResponse{
		Route:        "/login",
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusOK,
	}

	reqBodyStr, _ := json.Marshal(user)
	req := httptest.NewRequest(test.Method, test.Route, bytes.NewBufferString(string(reqBodyStr)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to login user controller")
	}

	var errResponse utils.ResponseError
	responseBodyBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(responseBodyBytes, &errResponse)

	defer testutils.DeleteUser(&email)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, errResponse.Error, utils.ErrorCode[utils.LoginInvalidPassword])
}

func TestLoginAccountNotVerified(t *testing.T) {
	email := "login_account_not_verified@example.com"
	testutils.CreateUser(&email, false)

	user := &models.ReqLoginUser{
		Email:    email,
		Password: testutils.StrPassword,
	}

	test := &testutils.TestResponse{
		Route:        "/login",
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusUnauthorized,
	}

	reqBodyStr, _ := json.Marshal(user)
	req := httptest.NewRequest(test.Method, test.Route, bytes.NewBufferString(string(reqBodyStr)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to login user controller")
	}

	var errResponse utils.ResponseError
	responseBodyBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(responseBodyBytes, &errResponse)

	defer testutils.DeleteUser(&email)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, errResponse.Error, utils.ErrorCode[utils.LoginAccountNotVerified])
}

func TestLoginSuccess(t *testing.T) {
	email := "login_success@example.com"
	testutils.CreateUser(&email, true)

	user := &models.ReqLoginUser{
		Email:    email,
		Password: testutils.StrPassword,
	}

	test := &testutils.TestResponse{
		Route:        "/login",
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusOK,
	}

	reqBodyStr, _ := json.Marshal(user)
	req := httptest.NewRequest(test.Method, test.Route, bytes.NewBufferString(string(reqBodyStr)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to login user controller")
	}

	testutils.DeleteUser(&email)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.NotEmpty(t, resp.Header.Get("Set-Cookie"))
}

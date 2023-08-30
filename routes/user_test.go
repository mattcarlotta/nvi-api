package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	resBody := testutils.ParseJSONErrorBody(&resp.Body)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.RegisterInvalidBody])
}

func TestRegisterUserInvalidBody(t *testing.T) {
	user := &models.ReqRegisterUser{
		Name: "invalid",
		// invalid email to trigger validation failure
		Email:    "invalidexample",
		Password: testutils.StrPassword,
	}

	test := &testutils.TestResponse{
		Route:        "/register",
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusBadRequest,
	}

	reqBody, _ := json.Marshal(user)
	req := httptest.NewRequest(test.Method, test.Route, bytes.NewBufferString(string(reqBody)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to register user controller")
	}

	resBody := testutils.ParseJSONErrorBody(&resp.Body)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.RegisterInvalidBody])
}

func TestRegisterEmailTaken(t *testing.T) {
	email := "taken_email@example.com"
	u, _ := testutils.CreateUser(email, false)

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

	reqBody, _ := json.Marshal(user)
	req := httptest.NewRequest(test.Method, test.Route, bytes.NewBufferString(string(reqBody)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to register user controller")
	}

	resBody := testutils.ParseJSONErrorBody(&resp.Body)

	defer testutils.DeleteUser(&u)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.RegisterEmailTaken])
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

	reqBody, _ := json.Marshal(user)
	req := httptest.NewRequest(test.Method, test.Route, bytes.NewBufferString(string(reqBody)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to register user controller")
	}

	resBody := testutils.ParseText(&resp.Body)

	defer testutils.RemoveUserByEmail(user.Email)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, resBody, fmt.Sprintf(
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

	resBody := testutils.ParseJSONErrorBody(&resp.Body)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.LoginInvalidBody])
}

func TestLoginUserInvalidBody(t *testing.T) {
	user := &models.ReqLoginUser{
		// invalid email to trigger validation failure
		Email:    "invalidexample",
		Password: testutils.StrPassword,
	}

	test := &testutils.TestResponse{
		Route:        "/login",
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusBadRequest,
	}

	reqBody, _ := json.Marshal(user)
	req := httptest.NewRequest(test.Method, test.Route, bytes.NewBufferString(string(reqBody)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to login user controller")
	}

	resBody := testutils.ParseJSONErrorBody(&resp.Body)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.LoginInvalidBody])
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

	reqBody, _ := json.Marshal(user)
	req := httptest.NewRequest(test.Method, test.Route, bytes.NewBufferString(string(reqBody)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to login user controller")
	}

	resBody := testutils.ParseJSONErrorBody(&resp.Body)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.LoginUnregisteredEmail])
}

func TestLoginInvalidPassword(t *testing.T) {
	email := "login_invalid_password@example.com"
	u, _ := testutils.CreateUser(email, false)

	badPassword := testutils.StrPassword + "4"
	user := &models.ReqLoginUser{
		Email:    email,
		Password: badPassword,
	}

	test := &testutils.TestResponse{
		Route:        "/login",
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusOK,
	}

	reqBody, _ := json.Marshal(user)
	req := httptest.NewRequest(test.Method, test.Route, bytes.NewBufferString(string(reqBody)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to login user controller")
	}

	resBody := testutils.ParseJSONErrorBody(&resp.Body)

	defer testutils.DeleteUser(&u)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.LoginInvalidPassword])
}

func TestLoginAccountNotVerified(t *testing.T) {
	email := "login_account_not_verified@example.com"
	u, _ := testutils.CreateUser(email, false)

	user := &models.ReqLoginUser{
		Email:    email,
		Password: testutils.StrPassword,
	}

	test := &testutils.TestResponse{
		Route:        "/login",
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusUnauthorized,
	}

	reqBody, _ := json.Marshal(user)
	req := httptest.NewRequest(test.Method, test.Route, bytes.NewBufferString(string(reqBody)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to login user controller")
	}

	resBody := testutils.ParseJSONErrorBody(&resp.Body)

	defer testutils.DeleteUser(&u)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.LoginAccountNotVerified])
}

func TestLoginSuccess(t *testing.T) {
	email := "login_success@example.com"
	u, _ := testutils.CreateUser(email, true)

	user := &models.ReqLoginUser{
		Email:    email,
		Password: testutils.StrPassword,
	}

	test := &testutils.TestResponse{
		Route:        "/login",
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusOK,
	}

	reqBody, _ := json.Marshal(user)
	req := httptest.NewRequest(test.Method, test.Route, bytes.NewBufferString(string(reqBody)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to login user controller")
	}

	testutils.DeleteUser(&u)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.NotEmpty(t, resp.Header.Get("Set-Cookie"))
}

func TestVerifyAccountInvalidToken(t *testing.T) {
	test := &testutils.TestResponse{
		Route:        "/verify/account",
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusUnauthorized,
	}

	req := httptest.NewRequest(test.Method, test.Route, nil)
	req.Header.Add("Content-Type", "text/plain; charset=utf-8")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to verify account user controller")
	}

	resBody := testutils.ParseJSONErrorBody(&resp.Body)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.VerifyAccountInvalidToken])
}

func TestVerifyAccountInvalidEmailToken(t *testing.T) {
	token, _, err := utils.GenerateUserToken("email_does_not_exist@example.com")
	if err != nil {
		log.Fatal("unable to generate a new user token")
	}
	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/verify/account?token=%s", token),
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusUnprocessableEntity,
	}

	req := httptest.NewRequest(test.Method, test.Route, nil)
	req.Header.Add("Content-Type", "text/plain; charset=utf-8")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to verify account user controller")
	}

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
}

func TestVerifyAccountEmailAlreadyVerified(t *testing.T) {
	u, token := testutils.CreateUser("already_verified@example.com", true)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/verify/account?token=%s", token),
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusNotModified,
	}

	req := httptest.NewRequest(test.Method, test.Route, nil)
	req.Header.Add("Content-Type", "text/plain; charset=utf-8")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to verify account user controller")
	}

	defer testutils.DeleteUser(&u)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
}

func TestVerifyAccountSuccess(t *testing.T) {
	email := "verify_account@example.com"
	u, token := testutils.CreateUser(email, false)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/verify/account?token=%s", token),
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusAccepted,
	}

	req := httptest.NewRequest(test.Method, test.Route, nil)
	req.Header.Add("Content-Type", "text/plain; charset=utf-8")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to verify account user controller")
	}

	resBody := testutils.ParseText(&resp.Body)

	defer testutils.DeleteUser(&u)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, resBody, fmt.Sprintf("Successfully verified %s!", email))
}

func TestResendAccountVerifyInvalidToken(t *testing.T) {
	test := &testutils.TestResponse{
		Route:        "/reverify/account",
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusBadRequest,
	}

	req := httptest.NewRequest(test.Method, test.Route, nil)
	req.Header.Add("Content-Type", "text/plain; charset=utf-8")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to reverify account user controller")
	}

	resBody := testutils.ParseJSONErrorBody(&resp.Body)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.ResendAccountVerificationInvalidEmail])
}

func TestResendAccountVerifyInvalidEmail(t *testing.T) {
	test := &testutils.TestResponse{
		Route:        "/reverify/account?email=not_a_register_user@example.com",
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusNotModified,
	}

	req := httptest.NewRequest(test.Method, test.Route, nil)
	req.Header.Add("Content-Type", "text/plain; charset=utf-8")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to reverify account user controller")
	}

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
}

func TestResendAccountVerifyEmailAlreadyVerified(t *testing.T) {
	email := "already_reverified@example.com"
	u, _ := testutils.CreateUser(email, true)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/reverify/account?email=%s", email),
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusNotModified,
	}

	req := httptest.NewRequest(test.Method, test.Route, nil)
	req.Header.Add("Content-Type", "text/plain; charset=utf-8")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to reverify account user controller")
	}

	defer testutils.DeleteUser(&u)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
}

func TestResendAccountVerifySuccess(t *testing.T) {
	email := "not_reverified@example.com"
	u, _ := testutils.CreateUser(email, false)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/reverify/account?email=%s", email),
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusAccepted,
	}

	req := httptest.NewRequest(test.Method, test.Route, nil)
	req.Header.Add("Content-Type", "text/plain; charset=utf-8")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to reverify account user controller")
	}

	resBody := testutils.ParseText(&resp.Body)

	defer testutils.DeleteUser(&u)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, resBody, fmt.Sprintf("Resent a verification email to %s.", email))
}

func TestSendResetPasswordInvalidEmail(t *testing.T) {
	test := &testutils.TestResponse{
		Route:        "/reset/password",
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusBadRequest,
	}

	req := httptest.NewRequest(test.Method, test.Route, nil)
	req.Header.Add("Content-Type", "text/plain; charset=utf-8")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to send reset password user controller")
	}

	resBody := testutils.ParseJSONErrorBody(&resp.Body)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.SendResetPasswordInvalidEmail])
}

func TestSendResetPasswordUnregisteredEmail(t *testing.T) {
	test := &testutils.TestResponse{
		Route:        "/reset/password?email=not_a_register_user@example.com",
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusNotModified,
	}

	req := httptest.NewRequest(test.Method, test.Route, nil)
	req.Header.Add("Content-Type", "text/plain; charset=utf-8")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to send reset password user controller")
	}

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
}

func TestSendResetPasswordSuccess(t *testing.T) {
	email := "reset_password@example.com"
	u, _ := testutils.CreateUser(email, false)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/reset/password?email=%s", email),
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusAccepted,
	}

	req := httptest.NewRequest(test.Method, test.Route, nil)
	req.Header.Add("Content-Type", "text/plain; charset=utf-8")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to send reset password user controller")
	}

	resBody := testutils.ParseText(&resp.Body)

	defer testutils.DeleteUser(&u)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, resBody, fmt.Sprintf("Sent a password reset email to %s.", email))
}

func TestUpdatePasswordEmptyBody(t *testing.T) {
	test := &testutils.TestResponse{
		Route:        "/update/password",
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusBadRequest,
	}

	req := httptest.NewRequest(test.Method, test.Route, nil)
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to update password user controller")
	}

	resBody := testutils.ParseJSONErrorBody(&resp.Body)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.UpdatePasswordInvalidBody])
}

func TestUpdatePasswordInvalidBody(t *testing.T) {
	user := &models.ReqUpdateUser{
		Password: testutils.StrPassword,
		Token:    "", // invalid token
	}

	test := &testutils.TestResponse{
		Route:        "/update/password",
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusBadRequest,
	}

	reqBody, _ := json.Marshal(user)
	req := httptest.NewRequest(test.Method, test.Route, bytes.NewBufferString(string(reqBody)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to update password user controller")
	}

	resBody := testutils.ParseJSONErrorBody(&resp.Body)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.UpdatePasswordInvalidBody])
}

func TestUpdatePasswordInvalidToken(t *testing.T) {
	user := &models.ReqUpdateUser{
		Password: testutils.StrPassword,
		Token:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.SAItwXkQ_rZnRZ52ZdmCAlkZUheO9cClR6EPhcd854E", // invalid token
	}

	test := &testutils.TestResponse{
		Route:        "/update/password",
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusUnauthorized,
	}

	reqBody, _ := json.Marshal(user)
	req := httptest.NewRequest(test.Method, test.Route, bytes.NewBufferString(string(reqBody)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to update password user controller")
	}

	resBody := testutils.ParseJSONErrorBody(&resp.Body)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.UpdatePasswordInvalidToken])
}

func TestUpdatePasswordSuccess(t *testing.T) {
	u, _ := testutils.CreateUser("update_password@example.com", true)

	user := &models.ReqUpdateUser{
		Password: testutils.StrPassword,
		Token:    string(*u.Token),
	}

	test := &testutils.TestResponse{
		Route:        "/update/password",
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusCreated,
	}

	reqBody, _ := json.Marshal(user)
	req := httptest.NewRequest(test.Method, test.Route, bytes.NewBufferString(string(reqBody)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to update password user controller")
	}

	resBody := testutils.ParseText(&resp.Body)

	defer testutils.DeleteUser(&u)

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	assert.Equal(t, resBody, "Your account has been updated with a new password!")
}

func GetAccountInfoSuccess(t *testing.T) {
	_, token := testutils.CreateUser("update_password@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/account",
		Method:       fiber.MethodGet,
		ExpectedCode: fiber.StatusOK,
	}

	req := httptest.NewRequest(test.Method, test.Route, nil)
	req.Header.Add("Content-Type", "text/plain; charset=utf-8")
	req.Header.Add("Cookie", fmt.Sprintf("SESSION_TOKEN=%s", token))

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to get account info user controller")
	}

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
}

func DeleteAccountSuccess(t *testing.T) {
	_, token := testutils.CreateUser("delete_account@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/delete/account",
		Method:       fiber.MethodDelete,
		ExpectedCode: fiber.StatusOK,
	}

	req := httptest.NewRequest(test.Method, test.Route, nil)
	req.Header.Add("Content-Type", "text/plain; charset=utf-8")
	req.Header.Add("Cookie", fmt.Sprintf("SESSION_TOKEN=%s", token))

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatal("failed to make request to get account info user controller")
	}

	assert.Equal(t, test.ExpectedCode, resp.StatusCode)
}

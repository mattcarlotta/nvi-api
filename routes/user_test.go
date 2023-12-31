package routes

import (
	"fmt"
	"log"
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

	req := testutils.CreateHTTPRequest(test)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer res.Body.Close()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
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

	req := testutils.CreateHTTPRequest(test, user)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer res.Body.Close()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.RegisterInvalidBody])
}

func TestRegisterEmailTaken(t *testing.T) {
	email := "taken_email@example.com"
	u, _, _ := testutils.CreateUser(email, false)

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

	req := testutils.CreateHTTPRequest(test, user)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
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

	req := testutils.CreateHTTPRequest(test, user)

	res := sendAppRequest(req)

	resBody := testutils.ParseText(&res.Body)

	defer func() {
		testutils.RemoveUserByEmail(user.Email)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody, fmt.Sprintf(
		"Welcome, %s! Please check your %s inbox for steps to verify your account.", user.Name, user.Email,
	))
}

func TestLoggedinSuccess(t *testing.T) {
	u, token, _ := testutils.CreateUser("loggedin_account@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/loggedin",
		Method:       fiber.MethodGet,
		ExpectedCode: fiber.StatusOK,
	}

	req := testutils.CreateAuthHTTPRequest(test, &token)

	res := sendAppRequest(req)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
}

func TestLoginUserEmptyBody(t *testing.T) {
	test := &testutils.TestResponse{
		Route:        "/login",
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusBadRequest,
	}

	req := testutils.CreateHTTPRequest(test)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer res.Body.Close()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
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

	req := testutils.CreateHTTPRequest(test, user)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer res.Body.Close()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
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

	req := testutils.CreateHTTPRequest(test, user)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer res.Body.Close()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.LoginUnregisteredEmail])
}

func TestLoginInvalidPassword(t *testing.T) {
	email := "login_invalid_password@example.com"
	u, _, _ := testutils.CreateUser(email, false)

	badPassword := testutils.StrPassword + "4"
	user := &models.ReqLoginUser{
		Email:    email,
		Password: badPassword,
	}

	test := &testutils.TestResponse{
		Route:        "/login",
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusUnauthorized,
	}

	req := testutils.CreateHTTPRequest(test, user)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.LoginInvalidPassword])
}

func TestLoginAccountNotVerified(t *testing.T) {
	email := "login_account_not_verified@example.com"
	u, _, _ := testutils.CreateUser(email, false)

	user := &models.ReqLoginUser{
		Email:    email,
		Password: testutils.StrPassword,
	}

	test := &testutils.TestResponse{
		Route:        "/login",
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusUnauthorized,
	}

	req := testutils.CreateHTTPRequest(test, user)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.LoginAccountNotVerified])
}

func TestLoginSuccess(t *testing.T) {
	email := "login_success@example.com"
	u, _, _ := testutils.CreateUser(email, true)

	user := &models.ReqLoginUser{
		Email:    email,
		Password: testutils.StrPassword,
	}

	test := &testutils.TestResponse{
		Route:        "/login",
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusOK,
	}

	req := testutils.CreateHTTPRequest(test, user)

	res := sendAppRequest(req)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.NotEmpty(t, res.Header.Get("Set-Cookie"))
}

func TestVerifyAccountInvalidToken(t *testing.T) {
	test := &testutils.TestResponse{
		Route:        "/verify/account",
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusUnauthorized,
	}

	req := testutils.CreateHTTPRequest(test)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer res.Body.Close()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
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

	req := testutils.CreateHTTPRequest(test)

	res := sendAppRequest(req)

	defer res.Body.Close()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
}

func TestVerifyAccountEmailAlreadyVerified(t *testing.T) {
	u, token, _ := testutils.CreateUser("already_verified@example.com", true)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/verify/account?token=%s", token),
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusOK,
	}

	req := testutils.CreateHTTPRequest(test)

	res := sendAppRequest(req)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
}

func TestVerifyAccountSuccess(t *testing.T) {
	email := "verify_account@example.com"
	u, _, authToken := testutils.CreateUser(email, false)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/verify/account?token=%s", authToken),
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusCreated,
	}

	req := testutils.CreateHTTPRequest(test)

	res := sendAppRequest(req)

	resBody := testutils.ParseText(&res.Body)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody, fmt.Sprintf("Successfully verified %s!", email))
}

func TestResendAccountVerifyInvalidToken(t *testing.T) {
	test := &testutils.TestResponse{
		Route:        "/reverify/account",
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusBadRequest,
	}

	req := testutils.CreateHTTPRequest(test)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer res.Body.Close()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.ResendAccountVerificationInvalidEmail])
}

func TestResendAccountVerifyInvalidEmail(t *testing.T) {
	test := &testutils.TestResponse{
		Route:        "/reverify/account?email=not_a_register_user@example.com",
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusOK,
	}

	req := testutils.CreateHTTPRequest(test)

	res := sendAppRequest(req)

	defer res.Body.Close()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
}

func TestResendAccountVerifyEmailAlreadyVerified(t *testing.T) {
	email := "already_reverified@example.com"
	u, _, _ := testutils.CreateUser(email, true)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/reverify/account?email=%s", email),
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusOK,
	}

	req := testutils.CreateHTTPRequest(test)

	res := sendAppRequest(req)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
}

func TestResendAccountVerifySuccess(t *testing.T) {
	email := "not_reverified@example.com"
	u, _, _ := testutils.CreateUser(email, false)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/reverify/account?email=%s", email),
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusCreated,
	}

	req := testutils.CreateHTTPRequest(test)

	res := sendAppRequest(req)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
}

func TestSendResetPasswordInvalidEmail(t *testing.T) {
	test := &testutils.TestResponse{
		Route:        "/reset/password",
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusBadRequest,
	}

	req := testutils.CreateHTTPRequest(test)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer res.Body.Close()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.SendResetPasswordInvalidEmail])
}

func TestSendResetPasswordUnregisteredEmail(t *testing.T) {
	test := &testutils.TestResponse{
		Route:        "/reset/password?email=not_a_register_user@example.com",
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusOK,
	}

	req := testutils.CreateHTTPRequest(test)

	res := sendAppRequest(req)

	defer res.Body.Close()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
}

func TestSendResetPasswordSuccess(t *testing.T) {
	email := "reset_password@example.com"
	u, _, _ := testutils.CreateUser(email, false)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/reset/password?email=%s", email),
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusCreated,
	}

	req := testutils.CreateHTTPRequest(test)

	res := sendAppRequest(req)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
}

func TestUpdatePasswordEmptyBody(t *testing.T) {
	test := &testutils.TestResponse{
		Route:        "/update/password",
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusBadRequest,
	}

	req := testutils.CreateHTTPRequest(test)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer res.Body.Close()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
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

	req := testutils.CreateHTTPRequest(test, user)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer res.Body.Close()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.UpdatePasswordInvalidBody])
}

func TestUpdatePasswordInvalidToken(t *testing.T) {
	user := &models.ReqUpdateUser{
		Password: testutils.StrPassword,
		// invalid token
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.SAItwXkQ_rZnRZ52ZdmCAlkZUheO9cClR6EPhcd854E",
	}

	test := &testutils.TestResponse{
		Route:        "/update/password",
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusUnauthorized,
	}

	req := testutils.CreateHTTPRequest(test, user)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer res.Body.Close()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.UpdatePasswordInvalidToken])
}

func TestUpdatePasswordSuccess(t *testing.T) {
	u, _, authToken := testutils.CreateUser("update_password@example.com", true)

	user := &models.ReqUpdateUser{
		Password: testutils.StrPassword,
		Token:    authToken,
	}

	test := &testutils.TestResponse{
		Route:        "/update/password",
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusCreated,
	}

	req := testutils.CreateHTTPRequest(test, user)

	res := sendAppRequest(req)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
}

func TestUpdateDisplayNameMissingName(t *testing.T) {
	u, token, _ := testutils.CreateUser("update_display_name_missing_name@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/update/name",
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusBadRequest,
	}

	req := testutils.CreateAuthHTTPRequest(test, &token)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.UpdateDisplayNameMissingName])
}

func TestUpdateDisplayNameSuccess(t *testing.T) {
	u, token, _ := testutils.CreateUser("update_display_name_success@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/update/name/?name=updated_name",
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusCreated,
	}

	req := testutils.CreateAuthHTTPRequest(test, &token)

	res := sendAppRequest(req)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
}

func TestUpdateAPIKeySuccess(t *testing.T) {
	u, token, _ := testutils.CreateUser("update_api_key_success@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/update/apikey",
		Method:       fiber.MethodPatch,
		ExpectedCode: fiber.StatusCreated,
	}

	req := testutils.CreateAuthHTTPRequest(test, &token)

	res := sendAppRequest(req)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
}

func TestGetAccountInfoSuccess(t *testing.T) {
	u, token, _ := testutils.CreateUser("get_account_info@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/account",
		Method:       fiber.MethodGet,
		ExpectedCode: fiber.StatusOK,
	}

	req := testutils.CreateAuthHTTPRequest(test, &token)

	res := sendAppRequest(req)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
}

func TestDeleteAccountSuccess(t *testing.T) {
	_, token, _ := testutils.CreateUser("delete_account@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/delete/account",
		Method:       fiber.MethodDelete,
		ExpectedCode: fiber.StatusOK,
	}

	req := testutils.CreateAuthHTTPRequest(test, &token)

	res := sendAppRequest(req)

	defer res.Body.Close()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
}

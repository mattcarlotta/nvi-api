package routes

import (
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mattcarlotta/nvi-api/models"
	"github.com/mattcarlotta/nvi-api/test"
	"github.com/mattcarlotta/nvi-api/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetSecretInvalidID(t *testing.T) {
	u, token := testutils.CreateUser("get_secret_invalid_id@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/secret/not_a_valid_secret_uuid",
		Method:       fiber.MethodGet,
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
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.GetSecretInvalidID])
}

func TestGetSecretNonExistentID(t *testing.T) {
	u, token := testutils.CreateUser("get_secret_non_existent_id@example.com", true)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/secret/%s", uuid.NewString()),
		Method:       fiber.MethodGet,
		ExpectedCode: fiber.StatusNotFound,
	}

	req := testutils.CreateAuthHTTPRequest(test, &token)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.GetSecretNonExistentID])
}

func TestGetSecretSuccess(t *testing.T) {
	u, token := testutils.CreateUser("get_secret_non_existent_id@example.com", true)
	_, s := testutils.CreateEnvironmentAndSecret("get_secret_env_success", "GET_SECRET_KEY", "env_value", token)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/secret/%s", s.ID.String()),
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

func TestGetSecretsByEnvironmentInvalidID(t *testing.T) {
	u, token := testutils.CreateUser("get_secrets_by_invalid_env_id@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/secrets/not_a_valid_env_uuid",
		Method:       fiber.MethodGet,
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
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.GetSecretsByEnvInvalidID])
}

func TestGetSecretsByEnvironmentNonExistentID(t *testing.T) {
	u, token := testutils.CreateUser("get_secrets_by_invalid_env_id@example.com", true)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/secrets/%s", uuid.NewString()),
		Method:       fiber.MethodGet,
		ExpectedCode: fiber.StatusNotFound,
	}

	req := testutils.CreateAuthHTTPRequest(test, &token)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.GetSecretsByEnvNonExistentID])
}

func TestGetSecretsByEnvironmentSuccess(t *testing.T) {
	u, token := testutils.CreateUser("get_secrets_by_env_success@example.com", true)
	e, _ := testutils.CreateEnvironmentAndSecret("get_secrets_by_env_success", "GET_SECRET_ENV_KEY", "env_value", token)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/secrets/%s", e.ID.String()),
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

func TestCreateSecretEmptyBody(t *testing.T) {
	u, token := testutils.CreateUser("create_secret_empty@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/create/secret",
		Method:       fiber.MethodPost,
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
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.CreateSecretInvalidBody])
}

func TestCreateSecretInvalidBody(t *testing.T) {
	u, token := testutils.CreateUser("secret_invalid_body@example.com", true)
	testutils.CreateEnvironment("create_env_invalid_body", token)

	secret := &models.ReqCreateSecret{
		EnvironmentIDs: []string{"not_valid_uuid"},
		Key:            "INVALID_ID",
		Value:          "invalid_id",
	}

	test := &testutils.TestResponse{
		Route:        "/create/secret",
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusBadRequest,
	}

	req := testutils.CreateAuthHTTPRequest(test, &token, secret)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.CreateSecretInvalidBody])
}

func TestCreateSecretNonExistentID(t *testing.T) {
	u, token := testutils.CreateUser("secret_non_existent_id@example.com", true)
	testutils.CreateEnvironment("create_secret_env_non_existent_id", token)

	secret := &models.ReqCreateSecret{
		EnvironmentIDs: []string{uuid.NewString()},
		Key:            "NON_EXISTENT_ID",
		Value:          "non_existent_id",
	}

	test := &testutils.TestResponse{
		Route:        "/create/secret",
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusNotFound,
	}

	req := testutils.CreateAuthHTTPRequest(test, &token, secret)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.CreateSecretNonExistentEnv])
}

func TestCreateSecretKeyAlreadyExists(t *testing.T) {
	u, token := testutils.CreateUser("secret_already_exists@example.com", true)
	e, _ := testutils.CreateEnvironmentAndSecret("secret_exists", "SECRET_EXISTS", "abc123", token)

	secret := &models.ReqCreateSecret{
		EnvironmentIDs: []string{e.ID.String()},
		Key:            "SECRET_EXISTS",
		Value:          "def456",
	}

	test := &testutils.TestResponse{
		Route:        "/create/secret",
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusConflict,
	}

	req := testutils.CreateAuthHTTPRequest(test, &token, secret)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.CreateSecretKeyAlreadyExists])
}

func TestCreateSecretSuccess(t *testing.T) {
	u, token := testutils.CreateUser("create_secret_success@example.com", true)
	e := testutils.CreateEnvironment("secret_success", token)

	key := "SECRET_TO_SUCCESS"
	secret := &models.ReqCreateSecret{
		EnvironmentIDs: []string{e.ID.String()},
		Key:            key,
		Value:          "never play the game",
	}

	test := &testutils.TestResponse{
		Route:        "/create/secret",
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusCreated,
	}

	req := testutils.CreateAuthHTTPRequest(test, &token, secret)

	res := sendAppRequest(req)

	resBody := testutils.ParseText(&res.Body)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody, fmt.Sprintf("Successfully created %s!", key))
}

func TestDeleteSecretInvalidID(t *testing.T) {
	u, token := testutils.CreateUser("delete_secret_empty@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/delete/secret/not_valid_uuid",
		Method:       fiber.MethodDelete,
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
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.DeleteSecretInvalidID])
}

func TestDeleteSecretNonExistentID(t *testing.T) {
	u, token := testutils.CreateUser("delete_secret_non_existent_id@example.com", true)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/delete/secret/%s", uuid.New()),
		Method:       fiber.MethodDelete,
		ExpectedCode: fiber.StatusNotFound,
	}

	req := testutils.CreateAuthHTTPRequest(test, &token)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.DeleteSecretNonExistentID])
}

func TestDeleteSecretSuccess(t *testing.T) {
	u, token := testutils.CreateUser("delete_secret_success@example.com", true)
	_, s := testutils.CreateEnvironmentAndSecret("delete_secret_success", "DELETE_SECRET", "abc123", token)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/delete/secret/%s", s.ID.String()),
		Method:       fiber.MethodDelete,
		ExpectedCode: fiber.StatusCreated,
	}

	req := testutils.CreateAuthHTTPRequest(test, &token)

	res := sendAppRequest(req)

	resBody := testutils.ParseText(&res.Body)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody, fmt.Sprintf("Successfully deleted the %s secret!", s.Key))
}

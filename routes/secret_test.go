package routes

import (
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	// "github.com/mattcarlotta/nvi-api/models"
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

	defer testutils.DeleteUser(&u)
	defer res.Body.Close()

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

	defer testutils.DeleteUser(&u)
	defer res.Body.Close()

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

	defer testutils.DeleteUser(&u)
	defer res.Body.Close()

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

	defer testutils.DeleteUser(&u)
	defer res.Body.Close()

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

	defer testutils.DeleteUser(&u)
	defer res.Body.Close()

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

	defer testutils.DeleteUser(&u)
	defer res.Body.Close()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
}

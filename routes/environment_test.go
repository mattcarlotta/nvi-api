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

func TestGetAllEnvironmentsSuccess(t *testing.T) {
	u, token := testutils.CreateUser("get_all_env@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/environments",
		Method:       fiber.MethodGet,
		ExpectedCode: fiber.StatusOK,
	}

	req := testutils.CreateAuthHttpRequest(test, &token)

	res := sendAppRequest(req)

	defer testutils.DeleteUser(&u)

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
}

func TestGetEnvironmentInvalidID(t *testing.T) {
	u, token := testutils.CreateUser("get_env_invalid_id@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/environment/not_a_uuid",
		Method:       fiber.MethodGet,
		ExpectedCode: fiber.StatusBadRequest,
	}

	req := testutils.CreateAuthHttpRequest(test, &token)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer testutils.DeleteUser(&u)

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.GetEnvironmentInvalidID])
}

func TestGetEnvironmentNonExistentID(t *testing.T) {
	u, token := testutils.CreateUser("get_env_non_existent_id@example.com", true)

	envUUID := uuid.New()

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/environment/%s", envUUID.String()),
		Method:       fiber.MethodGet,
		ExpectedCode: fiber.StatusNotFound,
	}

	req := testutils.CreateAuthHttpRequest(test, &token)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer testutils.DeleteUser(&u)

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.GetEnvironmentNonExistentID])
}

func TestGetEnvironmentSuccess(t *testing.T) {
	u, token := testutils.CreateUser("get_env_success_id@example.com", true)
	e := testutils.CreateEnvironment("env_success", token)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/environment/%s", e.ID.String()),
		Method:       fiber.MethodGet,
		ExpectedCode: fiber.StatusOK,
	}

	req := testutils.CreateAuthHttpRequest(test, &token)

	res := sendAppRequest(req)

	defer testutils.DeleteUser(&u)

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
}

func TestCurrentEnvironmentInvalidName(t *testing.T) {
	u, token := testutils.CreateUser("create_env_empty@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/create/environment/a%20b$#.",
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusBadRequest,
	}

	req := testutils.CreateAuthHttpRequest(test, &token)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer testutils.DeleteUser(&u)

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.CreateEnvironmentInvalidName])
}

func TestCurrentEnvironmentNameTaken(t *testing.T) {
	u, token := testutils.CreateUser("create_env_taken@example.com", true)
	testutils.CreateEnvironment("taken_environment_name", token)

	test := &testutils.TestResponse{
		Route:        "/create/environment/taken_environment_name",
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusConflict,
	}

	req := testutils.CreateAuthHttpRequest(test, &token)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer testutils.DeleteUser(&u)

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.CreateEnvironmentNameTaken])
}

func TestCurrentEnvironmentSuccess(t *testing.T) {
	u, token := testutils.CreateUser("create_new_env@example.com", true)

	envName := "this_is_a_new_env"
	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/create/environment/%s", envName),
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusCreated,
	}

	req := testutils.CreateAuthHttpRequest(test, &token)

	res := sendAppRequest(req)

	resBody := testutils.ParseText(&res.Body)

	defer testutils.DeleteUser(&u)

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody, fmt.Sprintf("Successfully created a(n) %s environment!", envName))
}

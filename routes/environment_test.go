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

func TestGetAllEnvironmentsSuccess(t *testing.T) {
	u, token := testutils.CreateUser("get_all_env@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/environments",
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

func TestGetEnvironmentInvalidID(t *testing.T) {
	u, token := testutils.CreateUser("get_env_invalid_id@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/environment/not_a_uuid",
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
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.GetEnvironmentInvalidID])
}

func TestGetEnvironmentNonExistentID(t *testing.T) {
	u, token := testutils.CreateUser("get_env_non_existent_id@example.com", true)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/environment/%s", uuid.NewString()),
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

	req := testutils.CreateAuthHTTPRequest(test, &token)

	res := sendAppRequest(req)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
}

// func TestCreateEnvironmentInvalidName(t *testing.T) {
// 	u, token := testutils.CreateUser("create_env_empty@example.com", true)

// 	test := &testutils.TestResponse{
// 		Route:        "/create/environment/a%20b$#.",
// 		Method:       fiber.MethodPost,
// 		ExpectedCode: fiber.StatusBadRequest,
// 	}

// 	req := testutils.CreateAuthHTTPRequest(test, &token)

// 	res := sendAppRequest(req)

// 	resBody := testutils.ParseJSONBodyError(&res.Body)

// 	defer func() {
// 		testutils.DeleteUser(&u)
// 		res.Body.Close()
// 	}()

// 	assert.Equal(t, test.ExpectedCode, res.StatusCode)
// 	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.CreateEnvironmentInvalidName])
// }

// func TestCreateEnvironmentNameTaken(t *testing.T) {
// 	u, token := testutils.CreateUser("create_env_taken@example.com", true)
// 	testutils.CreateEnvironment("taken_environment_name", token)

// 	test := &testutils.TestResponse{
// 		Route:        "/create/environment/taken_environment_name",
// 		Method:       fiber.MethodPost,
// 		ExpectedCode: fiber.StatusConflict,
// 	}

// 	req := testutils.CreateAuthHTTPRequest(test, &token)

// 	res := sendAppRequest(req)

// 	resBody := testutils.ParseJSONBodyError(&res.Body)

// 	defer func() {
// 		testutils.DeleteUser(&u)
// 		res.Body.Close()
// 	}()

// 	assert.Equal(t, test.ExpectedCode, res.StatusCode)
// 	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.CreateEnvironmentNameTaken])
// }

// func TestCreateEnvironmentSuccess(t *testing.T) {
// 	u, token := testutils.CreateUser("create_new_env@example.com", true)

// 	envName := "this_is_a_new_env"
// 	test := &testutils.TestResponse{
// 		Route:        fmt.Sprintf("/create/environment/%s", envName),
// 		Method:       fiber.MethodPost,
// 		ExpectedCode: fiber.StatusCreated,
// 	}

// 	req := testutils.CreateAuthHTTPRequest(test, &token)

// 	res := sendAppRequest(req)

// 	resBody := testutils.ParseText(&res.Body)

// 	defer func() {
// 		testutils.DeleteUser(&u)
// 		res.Body.Close()
// 	}()

// 	assert.Equal(t, test.ExpectedCode, res.StatusCode)
// 	assert.Equal(t, resBody, fmt.Sprintf("Successfully created a(n) %s environment!", envName))
// }

func TestDeleteEnvironmentInvalidID(t *testing.T) {
	u, token := testutils.CreateUser("delete_env_invalid_id@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/delete/environment/not_a_uuid",
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
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.DeleteEnvironmentInvalidID])
}

func TestDeleteEnvironmentNonExistentID(t *testing.T) {
	u, token := testutils.CreateUser("delete_env_non_existent_id@example.com", true)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/delete/environment/%s", uuid.NewString()),
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
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.DeleteEnvironmentNonExistentID])
}

func TestDeleteEnvironmentSuccess(t *testing.T) {
	u, token := testutils.CreateUser("delete_new_env@example.com", true)
	e := testutils.CreateEnvironment("delete_environment_name", token)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/delete/environment/%s", e.ID.String()),
		Method:       fiber.MethodDelete,
		ExpectedCode: fiber.StatusOK,
	}

	req := testutils.CreateAuthHTTPRequest(test, &token)

	res := sendAppRequest(req)

	resBody := testutils.ParseText(&res.Body)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody, fmt.Sprintf("Successfully deleted the %s environment!", e.Name))
}

func TestUpdateEnvironmentInvalidID(t *testing.T) {
	u, token := testutils.CreateUser("update_env_invalid_id@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/update/environment",
		Method:       fiber.MethodPut,
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
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.UpdateEnvironmentInvalidBody])
}

func TestUpdateEnvironmentNonExistentID(t *testing.T) {
	u, token := testutils.CreateUser("update_env_non_existent_id@example.com", true)

	env := &models.ReqUpdateEnv{
		ID:          uuid.NewString(),
		UpdatedName: "uuid_does_not_exist",
	}

	test := &testutils.TestResponse{
		Route:        "/update/environment",
		Method:       fiber.MethodPut,
		ExpectedCode: fiber.StatusNotFound,
	}

	req := testutils.CreateAuthHTTPRequest(test, &token, env)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.UpdateEnvironmentNonExistentID])
}

func TestUpdateEnvironmentSuccess(t *testing.T) {
	u, token := testutils.CreateUser("update_env_success@example.com", true)
	e := testutils.CreateEnvironment("update_env_success", token)

	updatedName := "updated_env_name"
	env := &models.ReqUpdateEnv{
		ID:          e.ID.String(),
		UpdatedName: updatedName,
	}

	test := &testutils.TestResponse{
		Route:        "/update/environment",
		Method:       fiber.MethodPut,
		ExpectedCode: fiber.StatusOK,
	}

	req := testutils.CreateAuthHTTPRequest(test, &token, env)

	res := sendAppRequest(req)

	resBody := testutils.ParseText(&res.Body)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody, fmt.Sprintf("Successfully updated the environment name to %s!", updatedName))
}

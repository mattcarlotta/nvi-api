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

func TestGetEnvironmentByIDInvalidID(t *testing.T) {
	u, token := testutils.CreateUser("get_env_invalid_id@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/environment/id/not_a_uuid",
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

func TestGetEnvironmentByIDNonExistentID(t *testing.T) {
	u, token := testutils.CreateUser("get_env_non_existent_id@example.com", true)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/environment/id/%s", uuid.NewString()),
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

func TestGetEnvironmentByIDSuccess(t *testing.T) {
	u, token := testutils.CreateUser("get_env_success_id@example.com", true)
	p := testutils.CreateProject("get_env_success_id_project", token)
	e := testutils.CreateEnvironment("get_env_success_id", p.ID, token)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/environment/id/%s", e.ID.String()),
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

func TestGetEnvironmentByNameInvalidName(t *testing.T) {
	u, token := testutils.CreateUser("get_env_invalid_name@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/environment/name/a@#bc%20",
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
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.GetEnvironmentInvalidName])
}

func TestGetEnvironmentByNameNonExistentName(t *testing.T) {
	u, token := testutils.CreateUser("get_env_non_existent_name@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/environment/name/non_existent_name",
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
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.GetEnvironmentNonExistentName])
}

func TestGetEnvironmentByNameSuccess(t *testing.T) {
	u, token := testutils.CreateUser("get_env_success_name@example.com", true)
	p := testutils.CreateProject("get_env_success_name_project", token)
	e := testutils.CreateEnvironment("get_env_success_name", p.ID, token)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/environment/name/%s", e.Name),
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

func TestCreateEnvironmentInvalidBody(t *testing.T) {
	u, token := testutils.CreateUser("create_env_empty@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/create/environment",
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
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.CreateEnvironmentInvalidBody])
}

func TestCreateEnvironmentInvalidProjectID(t *testing.T) {
	u, token := testutils.CreateUser("create_env_taken@example.com", true)

	env := &models.ReqCreateEnv{
		Name:      "taken_environment_name",
		ProjectID: uuid.NewString(),
	}

	test := &testutils.TestResponse{
		Route:        "/create/environment",
		Method:       fiber.MethodPost,
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
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.CreateEnvironmentInvalidProjectID])
}

func TestCreateEnvironmentNameTaken(t *testing.T) {
	u, token := testutils.CreateUser("create_env_taken@example.com", true)
	p := testutils.CreateProject("taken_environment_name_project", token)
	testutils.CreateEnvironment("taken_environment_name", p.ID, token)

	env := &models.ReqCreateEnv{
		Name:      "taken_environment_name",
		ProjectID: p.ID.String(),
	}

	test := &testutils.TestResponse{
		Route:        "/create/environment",
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusConflict,
	}

	req := testutils.CreateAuthHTTPRequest(test, &token, env)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.CreateEnvironmentNameTaken])
}

func TestCreateEnvironmentSuccess(t *testing.T) {
	u, token := testutils.CreateUser("create_new_env@example.com", true)
	p := testutils.CreateProject("create_new_env_project", token)

	env := &models.ReqCreateEnv{
		Name:      "create_new_env",
		ProjectID: p.ID.String(),
	}

	test := &testutils.TestResponse{
		Route:        "/create/environment",
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusCreated,
	}

	req := testutils.CreateAuthHTTPRequest(test, &token, env)

	res := sendAppRequest(req)

	resBody := testutils.ParseText(&res.Body)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody, fmt.Sprintf("Successfully created a(n) %s environment!", env.Name))
}

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
	p := testutils.CreateProject("delete_environment_name_project", token)
	e := testutils.CreateEnvironment("delete_environment_name", p.ID, token)

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
	p := testutils.CreateProject("update_env_success_project", token)
	e := testutils.CreateEnvironment("update_env_success", p.ID, token)

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

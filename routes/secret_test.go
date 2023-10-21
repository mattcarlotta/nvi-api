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

func TestGetSecretsByProjectAndEnvironmentMissingProjectName(t *testing.T) {
	u, token, _ := testutils.CreateUser("get_secrets_by_proj_env_missing_project_name@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/secrets/projectenvironment",
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
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.GetProjectInvalidName])
}

func TestGetSecretsByProjectAndEnvironmentMissingEnvironmentName(t *testing.T) {
	u, token, _ := testutils.CreateUser("get_secrets_by_proj_env_missing_env_name@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/secrets/projectenvironment?project=example_project",
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

func TestGetSecretsByProjectAndEnvironmentInvalidProjectName(t *testing.T) {
	u, token, _ := testutils.CreateUser("get_secrets_by_proj_env_invalid_project_name@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/secrets/projectenvironment?project=example_project&environment=example_environment",
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
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.GetProjectNonExistentName])
}

func TestGetSecretsByProjectAndEnvironmentInvalidEnvironmentName(t *testing.T) {
	u, token, _ := testutils.CreateUser("get_secrets_by_proj_env_invalid_env_name@example.com", true)
	p := testutils.CreateProject("get_secrets_by_proj_env_invalid_env_name", token)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/secrets/projectenvironment?project=%s&environment=example_environment", p.Name),
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

func TestGetSecretsByProjectAndEnvironmentSuccess(t *testing.T) {
	u, token, _ := testutils.CreateUser("get_secrets_by_proj_env_success@example.com", true)
	p, e, _ := testutils.CreateProjectAndEnvironmentAndSecret("get_secrets_by_proj_env_success", "proj_env_exists", "SECRET_EXISTS", "abc123", token)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/secrets/projectenvironment?project=%s&environment=%s", p.Name, e.Name),
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

func TestGetSecretInvalidID(t *testing.T) {
	u, token, _ := testutils.CreateUser("get_secret_invalid_id@example.com", true)

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
	u, token, _ := testutils.CreateUser("get_secret_non_existent_id@example.com", true)

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
	u, token, _ := testutils.CreateUser("get_secret_non_existent_id@example.com", true)
	_, _, s := testutils.CreateProjectAndEnvironmentAndSecret("get_secret_project", "get_secret_env_success", "GET_SECRET_KEY", "env_value", token)

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
	u, token, _ := testutils.CreateUser("get_secrets_by_invalid_env_id@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/secrets/id/not_a_valid_env_uuid",
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
	u, token, _ := testutils.CreateUser("get_secrets_by_invalid_env_id@example.com", true)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/secrets/id/%s", uuid.NewString()),
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
	u, token, _ := testutils.CreateUser("get_secrets_by_env_success@example.com", true)
	_, e, _ := testutils.CreateProjectAndEnvironmentAndSecret("get_secrets_by_env_project", "get_secrets_by_env_success", "GET_SECRET_ENV_KEY", "env_value", token)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/secrets/id/%s", e.ID.String()),
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

func TestSearchForSecretsByEnvironmentIDAndSecretKeyInvalidKey(t *testing.T) {
	u, token, _ := testutils.CreateUser("search_4_secrets_envId_secretKey_invalid_key@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/secrets/search",
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
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.SearchForSecretsByEnvAndSecretInvalidKey])
}

func TestSearchForSecretsByEnvironmentIDAndSecretKeyInvalidEnv(t *testing.T) {
	u, token, _ := testutils.CreateUser("search_4_secrets_envId_secretKey_invalid_env@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/secrets/search?key=ABC_123",
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

func TestSearchForSecretsByEnvironmentIDAndSecretKeySuccess(t *testing.T) {
	u, token, _ := testutils.CreateUser("search_4_secrets_envId_secretKey_success@example.com", true)
	_, e, s := testutils.CreateProjectAndEnvironmentAndSecret("search_4_secrets_envId_secretKey_success", "get_secrets_by_env_success", "GET_SECRET_ENV_KEY", "env_value", token)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/secrets/search?key=%s&environmentID=%s", s.Key, e.ID.String()),
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
	u, token, _ := testutils.CreateUser("create_secret_empty@example.com", true)

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
	u, token, _ := testutils.CreateUser("secret_invalid_body@example.com", true)

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

func TestCreateSecretNonExistentProject(t *testing.T) {
	u, token, _ := testutils.CreateUser("secret_non_existent_id@example.com", true)

	secret := &models.ReqCreateSecret{
		ProjectID:      uuid.NewString(),
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
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.CreateSecretNonExistentProject])
}

func TestCreateSecretNonExistentID(t *testing.T) {
	u, token, _ := testutils.CreateUser("secret_non_existent_id@example.com", true)
	p := testutils.CreateProject("secret_non_existent_id_project", token)

	secret := &models.ReqCreateSecret{
		ProjectID:      p.ID.String(),
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
	u, token, _ := testutils.CreateUser("secret_already_exists@example.com", true)
	p, e, _ := testutils.CreateProjectAndEnvironmentAndSecret("secret_exists_project", "secret_exists", "SECRET_EXISTS", "abc123", token)

	secret := &models.ReqCreateSecret{
		ProjectID:      p.ID.String(),
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
	u, token, _ := testutils.CreateUser("create_secret_success@example.com", true)
	p := testutils.CreateProject("secret_success_project", token)
	e := testutils.CreateEnvironment("secret_success", p.ID, token)

	key := "SECRET_TO_SUCCESS"
	secret := &models.ReqCreateSecret{
		ProjectID:      p.ID.String(),
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

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
}

func TestDeleteSecretInvalidID(t *testing.T) {
	u, token, _ := testutils.CreateUser("delete_secret_empty@example.com", true)

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
	u, token, _ := testutils.CreateUser("delete_secret_non_existent_id@example.com", true)

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
	u, token, _ := testutils.CreateUser("delete_secret_success@example.com", true)
	_, _, s := testutils.CreateProjectAndEnvironmentAndSecret("delete_secret_success_project", "delete_secret_success", "DELETE_SECRET", "abc123", token)

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
	assert.Equal(t, resBody, fmt.Sprintf("Successfully removed the %s secret!", s.Key))
}

func TestUpdateSecretEmptyBody(t *testing.T) {
	u, token, _ := testutils.CreateUser("update_secret_empty@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/update/secret",
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
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.UpdateSecretInvalidBody])
}

func TestUpdateSecretInvalidBody(t *testing.T) {
	u, token, _ := testutils.CreateUser("update_secret_invalid@example.com", true)

	secret := &models.ReqUpdateSecret{
		ID: uuid.NewString(),
		// invalid environment ids
		EnvironmentIDs: []string{},
		Key:            "UPDATE_KEY",
		Value:          "update",
	}

	test := &testutils.TestResponse{
		Route:        "/update/secret",
		Method:       fiber.MethodPut,
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
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.UpdateSecretInvalidBody])
}

func TestUpdateSecretNonExistentID(t *testing.T) {
	u, token, _ := testutils.CreateUser("update_secret_non_existent_id@example.com", true)

	secret := &models.ReqUpdateSecret{
		// non-existent secret uuid
		ID:             uuid.NewString(),
		EnvironmentIDs: []string{uuid.NewString()},
		Key:            "UPDATE_KEY",
		Value:          "update",
	}

	test := &testutils.TestResponse{
		Route:        "/update/secret",
		Method:       fiber.MethodPut,
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
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.UpdateSecretInvalidID])
}

func TestUpdateSecretNonExistentEnv(t *testing.T) {
	u, token, _ := testutils.CreateUser("update_secret_non_existent_env@example.com", true)
	_, _, s := testutils.CreateProjectAndEnvironmentAndSecret("update_secret_non_existent_project", "update_secret_non_existent_env", "UPDATE_SECRET", "abc123", token)

	secret := &models.ReqUpdateSecret{
		ID: s.ID.String(),
		// non-existent env uuid
		EnvironmentIDs: []string{uuid.NewString()},
		Key:            "UPDATE_KEY",
		Value:          "update",
	}

	test := &testutils.TestResponse{
		Route:        "/update/secret",
		Method:       fiber.MethodPut,
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
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.UpdateSecretNonExistentEnv])
}

func TestUpdateSecretKeyAlreadyExists(t *testing.T) {
	u, token, _ := testutils.CreateUser("update_secret_already_exists@example.com", true)
	p := testutils.CreateProject("update_secret_already_exists_project", token)
	e, _ := testutils.CreateEnvironmentAndSecret("env_1", p.ID, "TAKEN_KEY", "abc123", token)
	_, s := testutils.CreateEnvironmentAndSecret("env_2", p.ID, "NEW_KEY", "abc123", token)

	secret := &models.ReqUpdateSecret{
		ID:             s.ID.String(),
		EnvironmentIDs: []string{e.ID.String()},
		Key:            "TAKEN_KEY",
		Value:          "def456",
	}

	test := &testutils.TestResponse{
		Route:        "/update/secret",
		Method:       fiber.MethodPut,
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
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.UpdateSecretKeyAlreadyExists])
}

func TestUpdateSecretSuccess(t *testing.T) {
	u, token, _ := testutils.CreateUser("update_secret_success@example.com", true)
	_, e, s := testutils.CreateProjectAndEnvironmentAndSecret("update_secret_success_project", "update_secret_success", "UPDATE_KEY", "abc123", token)

	secret := &models.ReqUpdateSecret{
		ID:             s.ID.String(),
		EnvironmentIDs: []string{e.ID.String()},
		Key:            "UPDATED_KEY",
		Value:          "def456",
	}

	test := &testutils.TestResponse{
		Route:        "/update/secret",
		Method:       fiber.MethodPut,
		ExpectedCode: fiber.StatusOK,
	}

	req := testutils.CreateAuthHTTPRequest(test, &token, secret)

	res := sendAppRequest(req)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
}

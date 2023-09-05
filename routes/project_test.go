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

func TestGetAllProjectsSuccess(t *testing.T) {
	u, token := testutils.CreateUser("get_all_projects@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/projects",
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

func TestGetProjectByIDInvalidID(t *testing.T) {
	u, token := testutils.CreateUser("get_project_invalid_id@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/project/id/not_a_uuid",
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
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.GetProjectInvalidID])
}

func TestGetProjectByIDNonExistentID(t *testing.T) {
	u, token := testutils.CreateUser("get_project_non_existent_id@example.com", true)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/project/id/%s", uuid.NewString()),
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
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.GetProjectInvalidID])
}

func TestGetProjectByIDSuccess(t *testing.T) {
	u, token := testutils.CreateUser("get_project_success_id@example.com", true)
	p := testutils.CreateProject("get_project_success_id_project", token)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/project/id/%s", p.ID.String()),
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

func TestGetProjectByNameInvalidName(t *testing.T) {
	u, token := testutils.CreateUser("get_project_invalid_name@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/project/name/a@#bc%20",
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

func TestGetProjectByNameNonExistentName(t *testing.T) {
	u, token := testutils.CreateUser("get_project_non_existent_name@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/project/name/non_existent_name",
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

func TestGetProjectByNameSuccess(t *testing.T) {
	u, token := testutils.CreateUser("get_project_success_name@example.com", true)
	p := testutils.CreateProject("get_project_success_name_project", token)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/project/name/%s", p.Name),
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

func TestCreateProjectInvalidName(t *testing.T) {
	u, token := testutils.CreateUser("create_env_empty@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/create/project/a@#b%20c",
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
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.CreateProjectInvalidName])
}

func TestCreateProjectNameTaken(t *testing.T) {
	u, token := testutils.CreateUser("create_project_taken_name@example.com", true)
	p := testutils.CreateProject("create_project_taken_name", token)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/create/project/%s", p.Name),
		Method:       fiber.MethodPost,
		ExpectedCode: fiber.StatusConflict,
	}

	req := testutils.CreateAuthHTTPRequest(test, &token)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.CreateProjectNameTaken])
}

func TestCreateProjectSuccess(t *testing.T) {
	u, token := testutils.CreateUser("create_new_project@example.com", true)

	projectName := "create_new_project"
	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/create/project/%s", projectName),
		Method:       fiber.MethodPost,
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
	assert.Equal(t, resBody, fmt.Sprintf("Successfully created a(n) %s project!", projectName))
}

func TestDeleteProjectInvalidID(t *testing.T) {
	u, token := testutils.CreateUser("delete_project_invalid_id@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/delete/project/not_a_uuid",
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
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.DeleteProjectInvalidID])
}

func TestDeleteProjectNonExistentID(t *testing.T) {
	u, token := testutils.CreateUser("delete_project_non_existent_id@example.com", true)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/delete/project/%s", uuid.NewString()),
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
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.DeleteProjectNonExistentID])
}

func TestDeleteProjectSuccess(t *testing.T) {
	u, token := testutils.CreateUser("delete_project_id@example.com", true)
	p := testutils.CreateProject("delete_project_id", token)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/delete/project/%s", p.ID.String()),
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
	assert.Equal(t, resBody, fmt.Sprintf("Successfully deleted the %s project!", p.Name))
}

func TestUpdateProjectInvalidBody(t *testing.T) {
	u, token := testutils.CreateUser("update_project_invalid_body@example.com", true)

	test := &testutils.TestResponse{
		Route:        "/update/project",
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
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.UpdateProjectInvalidBody])
}

func TestUpdateProjectInvalidProjectID(t *testing.T) {
	u, token := testutils.CreateUser("update_project_invalid_project_id@example.com", true)

	project := &models.ReqUpdateProject{
		ID:          "not_a_uuid",
		UpdatedName: "invalid_project_id",
	}

	test := &testutils.TestResponse{
		Route:        "/update/project",
		Method:       fiber.MethodPut,
		ExpectedCode: fiber.StatusBadRequest,
	}

	req := testutils.CreateAuthHTTPRequest(test, &token, project)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.UpdateProjectInvalidBody])
}

func TestUpdateProjectNonExistentID(t *testing.T) {
	u, token := testutils.CreateUser("update_project_non_existent_id@example.com", true)

	project := &models.ReqUpdateProject{
		ID:          uuid.NewString(),
		UpdatedName: "project_uuid_does_not_exist",
	}

	test := &testutils.TestResponse{
		Route:        "/update/project",
		Method:       fiber.MethodPut,
		ExpectedCode: fiber.StatusNotFound,
	}

	req := testutils.CreateAuthHTTPRequest(test, &token, project)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.UpdateProjectNonExistentID])
}

func TestUpdateProjectNameTaken(t *testing.T) {
	u, token := testutils.CreateUser("update_project_name_taken@example.com", true)
	p1 := testutils.CreateProject("update_project_name_taken", token)
	p2 := testutils.CreateProject("update_project_name_unused", token)

	project := &models.ReqUpdateProject{
		ID:          p2.ID.String(),
		UpdatedName: p1.Name,
	}

	test := &testutils.TestResponse{
		Route:        "/update/project",
		Method:       fiber.MethodPut,
		ExpectedCode: fiber.StatusConflict,
	}

	req := testutils.CreateAuthHTTPRequest(test, &token, project)

	res := sendAppRequest(req)

	resBody := testutils.ParseJSONBodyError(&res.Body)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody.Error, utils.ErrorCode[utils.UpdateProjectNameTaken])
}

func TestUpdateProjectSuccess(t *testing.T) {
	u, token := testutils.CreateUser("update_project_success@example.com", true)
	p := testutils.CreateProject("update_project_success", token)

	updatedName := "updated_project_name_success"
	project := &models.ReqUpdateProject{
		ID:          p.ID.String(),
		UpdatedName: updatedName,
	}

	test := &testutils.TestResponse{
		Route:        "/update/project",
		Method:       fiber.MethodPut,
		ExpectedCode: fiber.StatusOK,
	}

	req := testutils.CreateAuthHTTPRequest(test, &token, project)

	res := sendAppRequest(req)

	resBody := testutils.ParseText(&res.Body)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody, fmt.Sprintf("Successfully updated the project name to %s!", updatedName))
}

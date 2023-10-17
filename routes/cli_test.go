package routes

import (
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/mattcarlotta/nvi-api/test"
	"github.com/stretchr/testify/assert"
)

func TestGetSecretByAPIKeyMissingAPIKey(t *testing.T) {
	test := &testutils.TestResponse{
		Route:        "/cli/secrets/",
		Method:       fiber.MethodGet,
		ExpectedCode: fiber.StatusUnauthorized,
	}

	req := testutils.CreateHTTPRequest(test)

	res := sendAppRequest(req)

	resBody := testutils.ParseText(&res.Body)

	defer func() {
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody, "a valid apiKey must be supplied in order to access secrets")
}

func TestGetSecretByAPIKeyInvalidAPIKey(t *testing.T) {
	test := &testutils.TestResponse{
		Route:        "/cli/secrets/?apiKey=notvalid",
		Method:       fiber.MethodGet,
		ExpectedCode: fiber.StatusNotFound,
	}

	req := testutils.CreateHTTPRequest(test)

	res := sendAppRequest(req)

	resBody := testutils.ParseText(&res.Body)

	defer func() {
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody, "the provided apiKey is not valid. please try again")
}

func TestGetSecretByAPIKeyMissingProject(t *testing.T) {
	u, _ := testutils.CreateUser("cli_get_secrets_missing_project@example.com", true)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/cli/secrets/?apiKey=%s", u.APIKey),
		Method:       fiber.MethodGet,
		ExpectedCode: fiber.StatusBadRequest,
	}

	req := testutils.CreateHTTPRequest(test)

	res := sendAppRequest(req)

	resBody := testutils.ParseText(&res.Body)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody, "a valid project name must be supplied in order to access secrets")
}

func TestGetSecretByAPIKeyInvalidProject(t *testing.T) {
	u, _ := testutils.CreateUser("cli_get_secrets_invalid_project@example.com", true)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/cli/secrets/?apiKey=%s&project=not_valid", u.APIKey),
		Method:       fiber.MethodGet,
		ExpectedCode: fiber.StatusNotFound,
	}

	req := testutils.CreateHTTPRequest(test)

	res := sendAppRequest(req)

	resBody := testutils.ParseText(&res.Body)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody, "unable to locate a project with the provided name")
}

func TestGetSecretByAPIKeyMissingEnvironment(t *testing.T) {
	u, token := testutils.CreateUser("cli_get_secrets_missing_enviroment@example.com", true)
	p := testutils.CreateProject("cli_get_secrets_missing_enviroment", token)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/cli/secrets/?apiKey=%s&project=%s", u.APIKey, p.Name),
		Method:       fiber.MethodGet,
		ExpectedCode: fiber.StatusBadRequest,
	}

	req := testutils.CreateHTTPRequest(test)

	res := sendAppRequest(req)

	resBody := testutils.ParseText(&res.Body)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody, "a valid environment name must be supplied in order to access secrets")
}

func TestGetSecretByAPIKeyInvalidEnvironment(t *testing.T) {
	u, token := testutils.CreateUser("cli_get_secrets_invalid_enviroment@example.com", true)
	p := testutils.CreateProject("cli_get_secrets_invalid_enviroment", token)

	envName := "not_valid"
	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/cli/secrets/?apiKey=%s&project=%s&environment=%s", u.APIKey, p.Name, envName),
		Method:       fiber.MethodGet,
		ExpectedCode: fiber.StatusNotFound,
	}

	req := testutils.CreateHTTPRequest(test)

	res := sendAppRequest(req)

	resBody := testutils.ParseText(&res.Body)

	defer func() {
		testutils.DeleteUser(&u)
		res.Body.Close()
	}()

	assert.Equal(t, test.ExpectedCode, res.StatusCode)
	assert.Equal(t, resBody, fmt.Sprintf("unable to locate a '%s' environment within the '%s' project", envName, p.Name))
}

func TestGetSecretByAPIKeySuccess(t *testing.T) {
	u, token := testutils.CreateUser("cli_get_secrets_success@example.com", true)
	p := testutils.CreateProject("cli_get_secrets_success", token)
	e, _ := testutils.CreateEnvironmentAndSecret("env_1", p.ID, "KEY", "abc123", token)

	test := &testutils.TestResponse{
		Route:        fmt.Sprintf("/cli/secrets/?apiKey=%s&project=%s&environment=%s", u.APIKey, p.Name, e.Name),
		Method:       fiber.MethodGet,
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

package controllers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
	"github.com/mattcarlotta/nvi-api/utils"
)

type ReqEnv struct {
	OriginalName string `json:"originalName"`
	UpdatedName  string `json:"updatedName"`
}

func CreateEnvironment(c *fiber.Ctx) error {
	var db = database.GetConnection()
	envName := c.Params("name")
	if len(envName) == 0 {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "You must provide a valid environment!")
	}

	var userSessionId = c.Locals("userSessionId").(string)

	var environment models.Environment
	if err := db.Where("name=? AND user_id=?", &envName, &userSessionId).First(&environment).Error; err == nil {
		return utils.SendErrorResponse(
			c,
			http.StatusOK,
			fmt.Sprintf("The provided environment '%s' already exists. Please choose a different environment name!", envName),
		)
	}

	newEnvironment := models.Environment{Name: envName, UserId: uuid.MustParse(userSessionId)}
	var err = db.Model(&environment).Create(&newEnvironment).Error
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusCreated).SendString(fmt.Sprintf("Successfully created the %s environment!", envName))
}

func DeleteEnvironment(c *fiber.Ctx) error {
	var db = database.GetConnection()
	id := c.Params("id")
	if len(id) == 0 {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "You must provide a valid environment!")
	}

	var userSessionId = c.Locals("userSessionId").(string)

	var environment models.Environment
	if err := db.Where("id=? AND user_id=?", &id, &userSessionId).First(&environment).Error; err != nil {
		return utils.SendErrorResponse(
			c,
			http.StatusOK,
			"The provided environment doesn't appear to exist!",
		)
	}

	var err = db.Delete(&environment).Error
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusCreated).SendString(fmt.Sprintf("Successfully deleted the %s environment!", environment.Name))
}

func UpdateEnvironment(c *fiber.Ctx) error {
	var db = database.GetConnection()
	data := new(ReqEnv)
	if err := c.BodyParser(data); err != nil {
		return utils.SendErrorResponse(
			c,
			http.StatusBadRequest,
			"You must provide a valid original environment name and an updated environment name!",
		)
	}

	// TODO(carlotta): Add field validations for "originalName" and "updatedName"
	if len(data.OriginalName) == 0 || len(data.UpdatedName) == 0 {
		return utils.SendErrorResponse(
			c,
			http.StatusBadRequest,
			"You must provide a valid original and updated name!",
		)
	}

	var userSessionId = c.Locals("userSessionId").(string)

	var environment models.Environment
	if err := db.Where("name=? AND user_id=?", &data.OriginalName, &userSessionId).First(&environment).Error; err != nil {
		return utils.SendErrorResponse(
			c,
			http.StatusOK,
			fmt.Sprintf("The provided environment '%s' doesn't appear to exist!", data.OriginalName),
		)
	}

	var err = db.Model(&environment).Update("name", &data.UpdatedName).Error
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusCreated).SendString(
		fmt.Sprintf("Successfully updated the environment from '%s' to '%s'!", data.OriginalName, data.UpdatedName),
	)
}

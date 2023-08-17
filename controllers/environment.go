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
	ID          string `json:"id"`
	UpdatedName string `json:"updatedName"`
}

func GetAllEnvironments(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionId := c.Locals("userSessionId").(string)

	var environments []models.Environment
	db.Where("user_id=?", &userSessionId).Find(&environments)

	return c.Status(http.StatusOK).JSON(environments)
}

func GetEnvironmentById(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionId := c.Locals("userSessionId").(string)

	id := c.Params("id")
	if len(id) == 0 {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "You must provide a valid environment id!")
	}

	var environment models.Environment
	if err := db.Where("id=? AND user_id=?", &id, &userSessionId).First(&environment).Error; err != nil {
		return utils.SendErrorResponse(
			c,
			http.StatusOK,
			"The provided environment doesn't appear to exist!",
		)
	}

	return c.Status(http.StatusOK).JSON(environment)
}

func CreateEnvironment(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionId := c.Locals("userSessionId").(string)

	envName := c.Params("name")
	if len(envName) == 0 {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "You must provide a valid environment name!")
	}

	var environment models.Environment
	if err := db.Where("name=? AND user_id=?", &envName, &userSessionId).First(&environment).Error; err == nil {
		return utils.SendErrorResponse(
			c,
			http.StatusOK,
			fmt.Sprintf(
				"The provided environment '%s' already exists. Please choose a different environment name!",
				envName,
			),
		)
	}

	newEnvironment := models.Environment{Name: envName, UserId: uuid.MustParse(userSessionId)}
	if err := db.Model(&environment).Create(&newEnvironment).Error; err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusCreated).SendString(fmt.Sprintf("Successfully created a(n) %s environment!", envName))
}

func DeleteEnvironment(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionId := c.Locals("userSessionId").(string)

	id := c.Params("id")
	if len(id) == 0 {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "You must provide a valid environment!")
	}

	var environment models.Environment
	if err := db.Where("id=? AND user_id=?", &id, &userSessionId).First(&environment).Error; err != nil {
		return utils.SendErrorResponse(
			c,
			http.StatusOK,
			"The provided environment doesn't appear to exist!",
		)
	}

	if err := db.Delete(&environment).Error; err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusCreated).SendString(
		fmt.Sprintf("Successfully deleted the %s environment!", environment.Name),
	)
}

func UpdateEnvironment(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionId := c.Locals("userSessionId").(string)

	data := new(ReqEnv)
	if err := c.BodyParser(data); err != nil {
		return utils.SendErrorResponse(
			c,
			http.StatusBadRequest,
			"You must provide a valid environment id an updated environment name!",
		)
	}

	// TODO(carlotta): Add field validations for "id" and "updatedName"
	if len(data.ID) == 0 || len(data.UpdatedName) == 0 {
		return utils.SendErrorResponse(
			c,
			http.StatusBadRequest,
			"You must provide a valid environment id an updated environment name!",
		)
	}

	var environment models.Environment
	if err := db.Where("id=? AND user_id=?", &data.ID, &userSessionId).First(&environment).Error; err != nil {
		return utils.SendErrorResponse(
			c,
			http.StatusOK,
			"The provided environment doesn't appear to exist.",
		)
	}

	if err := db.Model(&environment).Update("name", &data.UpdatedName).Error; err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusCreated).SendString(
		fmt.Sprintf("Successfully updated the environment name to %s!", data.UpdatedName),
	)
}

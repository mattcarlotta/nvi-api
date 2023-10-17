package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
	"github.com/mattcarlotta/nvi-api/utils"
	"gorm.io/gorm"
)

func GetSecretsByAPIKey(c *fiber.Ctx) error {
	db := database.GetConnection()
	apiKey := utils.GetAPIKey(c)

	return db.Transaction(func(tx *gorm.DB) error {
		var user models.User
		if err := tx.Where(&models.User{APIKey: apiKey}).First(&user).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString(
				"the provided apiKey is not valid. please try again",
			)
		}

		projectName := c.Query("project")
		if err := utils.Validate().Var(projectName, "required,name,lte=255"); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(
				"a valid project name must be supplied in order to access secrets",
			)
		}

		var project models.Project
		if err := tx.Where(
			&models.Project{Name: projectName, UserID: user.ID},
		).First(&project).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString(
				"unable to locate a project with the provided name",
			)
		}

		environmentName := c.Query("environment")
		if err := utils.Validate().Var(environmentName, "required,name,lte=255"); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(
				"a valid environment name must be supplied in order to access secrets",
			)
		}

		var environment models.Environment
		if err := tx.Where(
			&models.Environment{Name: environmentName, ProjectID: project.ID, UserID: user.ID},
		).First(&environment).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString(
				fmt.Sprintf("unable to locate a '%s' environment within the '%s' project", environmentName, projectName),
			)
		}

		var secrets []models.SecretResult
		if err := tx.Raw(
			utils.FindSecretsByEnvIDQuery, user.ID, utils.GenerateJSONIDString(environment.ID),
		).Scan(&secrets).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		var stringifiedSecrets string
		for _, secret := range secrets {
			decryptedValue, err := utils.DecryptSecretValue(secret.Value, secret.Nonce)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}
			var keyValue = secret.Key + "=" + string(decryptedValue) + "\n"
			stringifiedSecrets += keyValue
		}

		return c.Status(fiber.StatusOK).SendString(stringifiedSecrets)
	})
}

func GetProjectsByAPIKey(c *fiber.Ctx) error {
	db := database.GetConnection()
	apiKey := utils.GetAPIKey(c)

	var user models.User
	if err := db.Where(&models.User{APIKey: apiKey}).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString(
			"the provided apiKey is not valid. please try again",
		)
	}

	var projects []models.Project
	if err := db.Where(
		&models.Project{UserID: user.ID},
	).Find(&projects).Error; err != nil || len(projects) == 0 {
		return c.Status(fiber.StatusNotFound).SendString(
			"unable to locate any projects",
		)
	}

	var stringifiedProjects string
	for _, p := range projects {
		stringifiedProjects += p.Name + "\n"
	}

	return c.Status(fiber.StatusOK).SendString(stringifiedProjects)
}

func GetEnvironmentsByAPIKey(c *fiber.Ctx) error {
	db := database.GetConnection()
	apiKey := utils.GetAPIKey(c)

	return db.Transaction(func(tx *gorm.DB) error {
		var user models.User
		if err := tx.Where(&models.User{APIKey: apiKey}).First(&user).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString(
				"the provided apiKey is not valid. please try again",
			)
		}

		projectName := c.Query("project")
		if err := utils.Validate().Var(projectName, "required,name,lte=255"); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(
				"a valid project name must be supplied in order to access secrets",
			)
		}

		var project models.Project
		if err := tx.Where(
			&models.Project{Name: projectName, UserID: user.ID},
		).First(&project).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString(
				"unable to locate a project with the provided name",
			)
		}

		var environments []models.Environment
		if err := tx.Where(
			&models.Environment{ProjectID: project.ID, UserID: user.ID},
		).Find(&environments).Error; err != nil || len(environments) == 0 {
			return c.Status(fiber.StatusNotFound).SendString(
				fmt.Sprintf("unable to locate any environments within the '%s' project", projectName),
			)
		}

		var stringifiedEnvironments string
		for _, e := range environments {
			stringifiedEnvironments += e.Name + "\n"
		}

		return c.Status(fiber.StatusOK).SendString(stringifiedEnvironments)
	})
}

package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
	"github.com/mattcarlotta/nvi-api/utils"
	"gorm.io/gorm"
)

func GetEnvironmentByID(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	id := c.Params("id")
	if err := utils.Validate().Var(id, "required,uuid"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.GetEnvironmentInvalidID))
	}

	var environment models.Environment
	if err := db.Where(
		&models.Environment{ID: utils.MustParseUUID(id), UserID: userSessionID},
	).First(&environment).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.JSONError(utils.GetEnvironmentNonExistentID))
	}

	return c.Status(fiber.StatusOK).JSON(environment)
}

func GetAllEnvironmentByProjectName(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	projectName := c.Params("name")
	if err := utils.Validate().Var(projectName, "required,name,lte=255"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.GetProjectInvalidName))
	}

	var project models.Project
	if err := db.Where(
		&models.Project{Name: projectName, UserID: userSessionID},
	).First(&project).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.JSONError(utils.GetProjectNonExistentName))
	}

	var environments []models.Environment
	db.Where(&models.Environment{UserID: userSessionID, ProjectID: project.ID}).Find(&environments)

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"environments": environments,
			"project":      project,
		},
	)
}

func GetEnvironmentByNameAndProjectID(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	name := c.Query("name")
	if err := utils.Validate().Var(name, "required,name,lte=255"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.GetEnvironmentInvalidName))
	}

	projectID := c.Query("projectID")
	if err := utils.Validate().Var(projectID, "required,uuid"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.GetEnvironmentInvalidProjectID))
	}

	var environment models.Environment
	if err := db.Where(
		&models.Environment{Name: name, ProjectID: utils.MustParseUUID(projectID), UserID: userSessionID},
	).First(&environment).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.JSONError(utils.GetEnvironmentNonExistentName))
	}

	return c.Status(fiber.StatusOK).JSON(environment)
}

func SearchForEnvironmentsByNameAndProjectID(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	name := c.Query("name")
	if err := utils.Validate().Var(name, "required,name,lte=255"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.GetEnvironmentInvalidName))
	}

	projectID := c.Query("projectID")
	if err := utils.Validate().Var(projectID, "required,uuid"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.GetEnvironmentInvalidProjectID))
	}

	var environments []models.Environment
	if err := db.Where(
		"name ILIKE ? AND project_id=? AND user_id=?", "%"+name+"%", utils.MustParseUUID(projectID), userSessionID,
	).Find(&environments).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	return c.Status(fiber.StatusOK).JSON(environments)
}

func CreateEnvironment(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	var data models.ReqCreateEnv
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.CreateEnvironmentInvalidBody))
	}

	var project models.Project
	if err := db.Where(&models.Project{ID: utils.MustParseUUID(data.ProjectID)}).First(&project).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.JSONError(utils.CreateEnvironmentInvalidProjectID))
	}

	newEnv := models.Environment{Name: data.Name, ProjectID: project.ID, UserID: userSessionID}
	var environment models.Environment
	if err := db.Where(&newEnv).First(&environment).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(utils.JSONError(utils.CreateEnvironmentNameTaken))
	}

	// TODO(carlotta): add a limit to how many environments can be created per project and account

	if err := db.Create(&newEnv).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	return c.Status(fiber.StatusCreated).JSON(newEnv)
}

func DeleteEnvironment(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	id := c.Params("id")
	if err := utils.Validate().Var(id, "required,uuid"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.DeleteEnvironmentInvalidID))
	}

	return db.Transaction(func(tx *gorm.DB) error {
		parsedID := utils.MustParseUUID(id)

		var environment models.Environment
		if err := tx.Where(
			&models.Environment{ID: parsedID, UserID: userSessionID},
		).First(&environment).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(utils.JSONError(utils.DeleteEnvironmentNonExistentID))
		}

		var secrets []models.Secret
		if err := tx.Preload("Environments").Not("id", parsedID).Find(
			&secrets, "user_id=?", userSessionID,
		).Error; err != nil || len(secrets) > 0 {
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
			}

			for _, secret := range secrets {
				for _, env := range secret.Environments {
					if env.ID == environment.ID {
						tx.Delete(&secret)
					}
				}
			}
		}

		if err := tx.Delete(&environment).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
		}

		return c.Status(fiber.StatusOK).SendString(
			fmt.Sprintf("Successfully deleted the %s environment!", environment.Name),
		)
	})
}

func UpdateEnvironment(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	var data models.ReqUpdateEnv
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.UpdateEnvironmentInvalidBody))
	}

	if err := utils.Validate().Struct(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.UpdateEnvironmentInvalidBody))
	}

	return db.Transaction(func(tx *gorm.DB) error {
		projectID := utils.MustParseUUID(data.ProjectID)

		var project models.Project
		if err := tx.Where(
			&models.Project{ID: projectID, UserID: userSessionID},
		).First(&project).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(utils.JSONError(utils.UpdateEnvironmentInvalidProjectID))
		}

		envID := utils.MustParseUUID(data.ID)

		if err := tx.Not(
			"id", envID,
		).Where(
			&models.Environment{Name: data.UpdatedName, ProjectID: project.ID, UserID: userSessionID},
		).First(&models.Environment{}).Error; err == nil {
			return c.Status(fiber.StatusConflict).JSON(utils.JSONError(utils.UpdateEnvironmentNameTaken))
		}

		var environment models.Environment
		if err := tx.Where(
			&models.Environment{ID: envID, ProjectID: projectID, UserID: userSessionID},
		).First(&environment).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(utils.JSONError(utils.UpdateEnvironmentNonExistentID))
		}

		if err := tx.Model(&environment).Update("name", data.UpdatedName).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
		}

		return c.Status(fiber.StatusOK).SendString(
			fmt.Sprintf("Successfully updated the environment name to %s!", data.UpdatedName),
		)
	})
}

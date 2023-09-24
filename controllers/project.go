package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
	"github.com/mattcarlotta/nvi-api/utils"
	"gorm.io/gorm"
)

func GetAllProjects(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	var projects []models.Project
	db.Where(&models.Project{UserID: userSessionID}).Find(&projects)

	return c.Status(fiber.StatusOK).JSON(projects)
}

func GetProjectByID(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	id := c.Params("id")
	if err := utils.Validate().Var(id, "required,uuid"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.GetProjectInvalidID))
	}

	var project models.Project
	if err := db.Where(
		&models.Project{ID: utils.MustParseUUID(id), UserID: userSessionID},
	).First(&project).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.JSONError(utils.GetProjectInvalidID))
	}

	return c.Status(fiber.StatusOK).JSON(project)
}

func GetProjectByName(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	name := c.Params("name")
	if err := utils.Validate().Var(name, "required,name,lte=255"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.GetProjectInvalidName))
	}

	var project models.Project
	if err := db.Where(
		&models.Project{Name: name, UserID: userSessionID},
	).First(&project).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.JSONError(utils.GetProjectNonExistentName))
	}

	return c.Status(fiber.StatusOK).JSON(project)
}

func SearchForProjectsByName(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	name := c.Params("name")
	if err := utils.Validate().Var(name, "required,name,lte=255"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.GetProjectInvalidName))
	}

	var projects []models.Project
	if err := db.Where(
		"name ILIKE ? AND user_id=?", "%"+name+"%", userSessionID,
	).Find(&projects).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	return c.Status(fiber.StatusOK).JSON(projects)
}

func CreateProject(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	name := c.Params("name")
	if err := utils.Validate().Var(name, "required,name,lte=255"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.CreateProjectInvalidName))
	}

	var project models.Project
	if err := db.Where(
		&models.Project{Name: name, UserID: userSessionID},
	).First(&project).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(utils.JSONError(utils.CreateProjectNameTaken))
	}

	// TODO(carlotta): add a limit to how many projects can be created per project and account
	newProject := models.Project{Name: name, UserID: userSessionID}
	if err := db.Create(&newProject).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	return c.Status(fiber.StatusCreated).SendString(fmt.Sprintf("Successfully created a(n) %s project!", newProject.Name))
}

func DeleteProject(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	id := c.Params("id")
	if err := utils.Validate().Var(id, "required,uuid"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.DeleteProjectInvalidID))
	}

	var project models.Project
	if err := db.Where(
		&models.Project{ID: utils.MustParseUUID(id), UserID: userSessionID},
	).First(&project).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.JSONError(utils.DeleteProjectNonExistentID))
	}

	if err := db.Delete(&project).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	return c.Status(fiber.StatusOK).SendString(
		fmt.Sprintf("Successfully removed the %s project!", project.Name),
	)
}

func UpdateProject(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	var data models.ReqUpdateProject
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.UpdateProjectInvalidBody))
	}

	if err := utils.Validate().Struct(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.UpdateProjectInvalidBody))
	}

	return db.Transaction(func(tx *gorm.DB) error {
		projectID := utils.MustParseUUID(data.ID)

		var existingProject models.Project
		if err := tx.Where(
			&models.Project{ID: projectID, UserID: userSessionID},
		).First(&existingProject).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(utils.JSONError(utils.UpdateProjectNonExistentID))
		}

		if err := tx.Not(
			"id", projectID,
		).Where(
			&models.Project{Name: data.UpdatedName, UserID: userSessionID},
		).First(&models.Project{}).Error; err == nil {
			return c.Status(fiber.StatusConflict).JSON(utils.JSONError(utils.UpdateProjectNameTaken))
		}

		if err := db.Model(&existingProject).Update("name", data.UpdatedName).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
		}

		return c.Status(fiber.StatusOK).SendString(
			fmt.Sprintf("Successfully updated the project name to %s!", data.UpdatedName),
		)
	})
}

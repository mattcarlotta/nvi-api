package controllers

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
	"github.com/mattcarlotta/nvi-api/utils"
)

func Register(c *fiber.Ctx) error {
	db := database.GetConnection()

	var data models.ReqRegisterUser
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.RegisterInvalidBody))
	}

	if err := utils.Validate().Struct(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.RegisterInvalidBody))
	}

	var user models.User
	if err := db.Where("email=?", data.Email).First(&user).Error; err == nil {
		return c.Status(fiber.StatusOK).JSON(utils.JSONError(utils.RegisterEmailTaken))
	}

	token, _, err := utils.GenerateUserToken(data.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	newToken := []byte(token)
	if err = db.Create(
		&models.User{
			Email:    data.Email,
			Name:     data.Name,
			Password: []byte(data.Password),
			Token:    &newToken,
		},
	).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	// TODO(carlotta): Send account verification email

	return c.Status(fiber.StatusCreated).SendString(
		fmt.Sprintf("Welcome, %s! Please check your %s inbox for steps to verify your account.", data.Name, data.Email),
	)
}

func Loggedin(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	var loggedInUser models.User
	if err := db.Where(&models.User{ID: userSessionID}).Omit("password").First(&loggedInUser).Error; err != nil {
		newError := errors.New(
			"encountered an unexpected error. Unable to locate the associated account from the current session",
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(newError))
	}

	return c.Status(fiber.StatusOK).JSON(loggedInUser)
}

func Login(c *fiber.Ctx) error {
	db := database.GetConnection()

	var data models.ReqLoginUser
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.LoginInvalidBody))
	}

	if err := utils.Validate().Struct(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.LoginInvalidBody))
	}

	var existingUser models.User
	if err := db.Where(&models.User{Email: data.Email}).First(&existingUser).Error; err != nil {
		return c.Status(fiber.StatusOK).JSON(utils.JSONError(utils.LoginUnregisteredEmail))
	}

	if !existingUser.MatchPassword(data.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.JSONError(utils.LoginInvalidPassword))
	}

	if !existingUser.Verified {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.JSONError(utils.LoginAccountNotVerified))
	}

	token, exp, err := existingUser.GenerateSessionToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	utils.SetSessionCookie(c, token, exp)
	c.Status(fiber.StatusOK)
	return nil
}

func Logout(c *fiber.Ctx) error {
	utils.SetSessionCookie(c, "", time.Unix(0, 0))
	c.Status(fiber.StatusOK)
	return nil
}

func VerifyAccount(c *fiber.Ctx) error {
	db := database.GetConnection()

	parsedToken, err := utils.ValidateUserToken(c.Query("token"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.JSONError(utils.VerifyAccountInvalidToken))
	}

	var user models.User
	if err := db.Where(&models.User{Email: parsedToken.Email}).First(&user).Error; err != nil {
		c.Status(fiber.StatusUnprocessableEntity)
		return nil
	}

	if user.Verified {
		c.Status(fiber.StatusNotModified)
		return nil
	}

	var newToken []byte
	if err = db.Model(&user).Updates(&models.User{Verified: true, Token: &newToken}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	return c.Status(fiber.StatusAccepted).SendString(fmt.Sprintf("Successfully verified %s!", user.Email))
}

func ResendAccountVerification(c *fiber.Ctx) error {
	db := database.GetConnection()

	email := c.Query("email")
	if err := utils.Validate().Var(email, "required,email,lte=255"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.ResendAccountVerificationInvalidEmail))
	}

	var user models.User
	if err := db.Where(&models.User{Email: email}).First(&user).Error; err != nil {
		c.Status(fiber.StatusNotModified)
		return nil
	}

	if user.Verified {
		c.Status(fiber.StatusNotModified)
		return nil
	}

	token, _, err := utils.GenerateUserToken(email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if err = db.Model(&user).Update("token", &token).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// TODO(carlotta): Send account verification email

	return c.Status(fiber.StatusAccepted).SendString(fmt.Sprintf("Resent a verification email to %s.", user.Email))
}

func SendResetPasswordEmail(c *fiber.Ctx) error {
	db := database.GetConnection()

	email := c.Query("email")
	if err := utils.Validate().Var(email, "required,email,lte=255"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.SendResetPasswordInvalidEmail))
	}

	var user models.User
	if err := db.Where(&models.User{Email: email}).First(&user).Error; err != nil {
		c.Status(fiber.StatusNotModified)
		return nil
	}

	token, _, err := utils.GenerateUserToken(email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if err = db.Model(&user).Update("token", &token).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// TODO(carlotta): Send account verification email

	return c.Status(fiber.StatusAccepted).SendString(fmt.Sprintf("Sent a password reset email to %s.", user.Email))
}

func UpdatePassword(c *fiber.Ctx) error {
	db := database.GetConnection()

	var data models.ReqUpdateUser
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.UpdatePasswordInvalidBody))
	}

	if err := utils.Validate().Struct(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.UpdatePasswordInvalidBody))
	}

	parsedToken, err := utils.ValidateUserToken(data.Token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.JSONError(utils.UpdatePasswordInvalidToken))
	}

	var user models.User
	if err := db.Where("email=? AND token IS NOT NULL", &parsedToken.Email).First(&user).Error; err != nil {
		c.Status(fiber.StatusNotModified)
		return nil
	}

	newPassword, err := utils.CreateEncryptedText([]byte(data.Password))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var newToken []byte
	if err = db.Model(&user).Updates(models.User{Password: newPassword, Token: &newToken}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).SendString("Your account has been updated with a new password!")
}

func GetAccountInfo(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	var user models.User
	if err := db.Where(&models.User{ID: userSessionID}).First(&user).Error; err != nil {
		newError := errors.New(
			"encountered an unexpected error. Unable to locate the associated account from the current session",
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(newError))
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func DeleteAccount(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	var user models.User
	if err := db.Where(&models.User{ID: userSessionID}).First(&user).Error; err != nil {
		newError := errors.New(
			"encountered an unexpected error. Unable to locate the associated account from the current session",
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(newError))
	}

	if err := db.Delete(&user).Error; err != nil {
		newError := errors.New(
			"encountered an unexpected error. Unable to delete user account",
		)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(newError))
	}

	return Logout(c)
}

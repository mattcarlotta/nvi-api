package controllers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
	"github.com/mattcarlotta/nvi-api/utils"
)

func Register(c *fiber.Ctx) error {
	db := database.GetConnection()

	data := new(models.ReqRegisterUser)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.RegisterEmptyBody))
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

func Login(c *fiber.Ctx) error {
	db := database.GetConnection()

	data := new(models.ReqLoginUser)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.LoginEmptyBody))
	}

	if err := utils.Validate().Struct(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.LoginInvalidBody))
	}

	var existingUser models.User
	if err := db.Where(&models.User{Email: data.Email}).First(&existingUser).Error; err != nil {
		return c.Status(fiber.StatusOK).JSON(utils.JSONError(utils.LoginUnregisteredEmail))
	}

	if !existingUser.MatchPassword(data.Password) {
		return c.Status(fiber.StatusOK).JSON(utils.JSONError(utils.LoginInvalidPassword))
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
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{"error": "The provided token is not valid. If the account verification token was sent over 30 " +
				"days ago, you will need to generate another account verification email.",
			},
		)
	}

	var user models.User
	if err := db.Where(&models.User{Email: parsedToken.Email}).First(&user).Error; err != nil {
		c.Status(fiber.StatusNotModified)
		return nil
	}

	if user.Verified {
		c.Status(fiber.StatusNotModified)
		return nil
	}

	var newToken []byte
	if err = db.Model(&user).Updates(&models.User{Verified: true, Token: &newToken}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("Successfully verified %s!", user.Email))
}

func ResendAccountVerification(c *fiber.Ctx) error {
	db := database.GetConnection()

	data := new(models.ReqEmailUser)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid email!"},
		)
	}

	if err := utils.Validate().Struct(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid email address to resend an account verification email to."},
		)
	}

	var user models.User
	if err := db.Where(&models.User{Email: data.Email}).First(&user).Error; err != nil {
		c.Status(fiber.StatusNotModified)
		return nil
	}

	if user.Verified {
		c.Status(fiber.StatusNotModified)
		return nil
	}

	token, _, err := utils.GenerateUserToken(data.Email)
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

	data := new(models.ReqEmailUser)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid email!"},
		)
	}

	if err := utils.Validate().Struct(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid email address to send an password email to."},
		)
	}

	var user models.User
	if err := db.Where(&models.User{Email: data.Email}).First(&user).Error; err != nil {
		c.Status(fiber.StatusNotModified)
		return nil
	}

	token, _, err := utils.GenerateUserToken(data.Email)
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

	data := new(models.ReqUpdateUser)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid password and password reset token."},
		)
	}

	if err := utils.Validate().Struct(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid password and password reset token."},
		)
	}

	parsedToken, err := utils.ValidateUserToken(data.Token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{"error": "The provided token is not valid. If the token was sent over 30 days ago, you will " +
				"need to generate another reset password email.",
			},
		)
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
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Encountered an unexpected error. Unable to locate the associated account."},
		)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func DeleteAccount(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	var user models.User
	if err := db.Where(&models.User{ID: userSessionID}).First(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Encountered an unexpected error. Unable to locate the associated account."},
		)
	}

	if err := db.Delete(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Encountered an unexpected error. Unable to delete account."},
		)
	}

	return Logout(c)
}

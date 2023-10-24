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

	newUser := models.User{
		Email:    data.Email,
		Name:     data.Name,
		Password: []byte(data.Password),
	}
	if err = db.Create(&newUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	if err = utils.SendAccountVerificationEmail(newUser.Name, newUser.Email, token); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	return c.Status(fiber.StatusCreated).SendString(
		fmt.Sprintf("Welcome, %s! Please check your %s inbox for steps to verify your account.", data.Name, data.Email),
	)
}

func Loggedin(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	var loggedInUser models.User
	if err := db.Where(&models.User{ID: userSessionID}).Omit("password").First(&loggedInUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(
			errors.New("unable to locate the associated account from the current session")),
		)
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
		c.Status(fiber.StatusOK)
		return nil
	}

	if err = db.Model(&user).Updates(&models.User{Verified: true}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	return c.Status(fiber.StatusCreated).SendString(fmt.Sprintf("Successfully verified %s!", user.Email))
}

func ResendAccountVerification(c *fiber.Ctx) error {
	db := database.GetConnection()

	email := c.Query("email")
	if err := utils.Validate().Var(email, "required,email,lte=255"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.ResendAccountVerificationInvalidEmail))
	}

	var user models.User
	if err := db.Where(&models.User{Email: email}).First(&user).Error; err != nil {
		c.Status(fiber.StatusOK)
		return nil
	}

	if user.Verified {
		c.Status(fiber.StatusOK)
		return nil
	}

	token, _, err := utils.GenerateUserToken(email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	if err = utils.SendAccountVerificationEmail(user.Name, user.Email, token); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	c.Status(fiber.StatusCreated)
	return nil
}

func SendResetPasswordEmail(c *fiber.Ctx) error {
	db := database.GetConnection()

	email := c.Query("email")
	if err := utils.Validate().Var(email, "required,email,lte=255"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.SendResetPasswordInvalidEmail))
	}

	var user models.User
	if err := db.Where(&models.User{Email: email}).First(&user).Error; err != nil {
		c.Status(fiber.StatusOK)
		return nil
	}

	token, _, err := utils.GenerateUserToken(email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	if err = utils.SendPasswordResetEmail(user.Name, user.Email, token); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	c.Status(fiber.StatusCreated)
	return nil
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
	if err := db.Where(&models.User{Email: parsedToken.Subject}).First(&user).Error; err != nil {
		c.Status(fiber.StatusOK)
		return nil
	}

	newPassword, err := utils.CreateEncryptedText([]byte(data.Password))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	if err = db.Model(&user).Update("password", newPassword).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	if err = utils.SendPasswordResetConfirmationEmail(user.Name, user.Email); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	c.Status(fiber.StatusCreated)
	return nil
}

func UpdateDisplayName(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	name := c.Query("name")
	if err := utils.Validate().Var(name, "required,gte=2,lte=64"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.JSONError(utils.UpdateDisplayNameMissingName))
	}

	var existingUser models.User
	if err := db.Where(&models.User{ID: userSessionID}).First(&existingUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	if err := db.Model(&existingUser).Update("name", name).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	token, exp, err := existingUser.GenerateSessionToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	utils.SetSessionCookie(c, token, exp)
	c.Status(fiber.StatusCreated)
	return nil
}

func UpdateAPIKey(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	var user models.User
	if err := db.Where(&models.User{ID: userSessionID}).First(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(
			errors.New("unable to locate the associated account from the current session")),
		)
	}

	newAPIKey := utils.CreateBase64EncodedUUID()
	if err := db.Model(&user).Update("APIKey", newAPIKey).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(err))
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"apiKey": newAPIKey})
}

func GetAccountInfo(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	var user models.User
	if err := db.Where(&models.User{ID: userSessionID}).First(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(
			errors.New("unable to locate the associated account from the current session")),
		)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func DeleteAccount(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionID := utils.GetSessionID(c)

	var user models.User
	if err := db.Where(&models.User{ID: userSessionID}).First(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(
			errors.New("unable to locate the associated account from the current session")),
		)
	}

	if err := db.Delete(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.UnknownJSONError(
			errors.New("unable to delete user account")),
		)
	}

	// TODO(carlotta): send out account deletion confirmation email

	return Logout(c)
}

package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
	"github.com/mattcarlotta/nvi-api/utils"
)

type ReqUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

func Register(c *fiber.Ctx) error {
	var db = database.GetConnection()

	data := new(ReqUser)
	if err := c.BodyParser(data); err != nil {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "You must provide a valid name, email and password!")
	}

	// TODO(carlotta): Add field validations for "name," "email," and "password"
	if len(data.Email) == 0 || len(data.Password) == 0 || len(data.Name) == 0 {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "You must provide a valid name, email and password!")
	}

	var user models.User
	if err := db.Where("email=?", &data.Email).First(&user).Error; err == nil {
		return utils.SendErrorResponse(
			c,
			http.StatusOK,
			fmt.Sprintf("The provided email '%s' may already exist or is not a valid email address.", data.Email),
		)
	}

	token, _, err := utils.GenerateUserToken(data.Email)
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	var newToken = []byte(token)
	newUser := models.User{
		Email:    data.Email,
		Name:     data.Name,
		Password: []byte(data.Password),
		Token:    &newToken,
	}

	err = db.Model(&user).Create(&newUser).Error
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	// TODO(carlotta): Send account verification email

	return c.Status(http.StatusCreated).SendString(
		fmt.Sprintf("Welcome, %s! Please check your %s inbox for steps to verify your account.", data.Name, data.Email),
	)
}

func Login(c *fiber.Ctx) error {
	var db = database.GetConnection()

	data := new(ReqUser)
	if err := c.BodyParser(data); err != nil {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "You must provide a valid name, email and password!")
	}

	// TODO(carlotta): Add field validations for "email" and "password"
	if len(data.Email) == 0 || len(data.Password) == 0 {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "You must provide a valid email and password!")
	}

	var existingUser models.User
	if err := db.Where("email=?", &data.Email).First(&existingUser).Error; err != nil {
		return utils.SendErrorResponse(
			c,
			http.StatusNoContent,
			fmt.Sprintf("The provided email '%s' may not exist or the provided password is incorrect.", data.Email),
		)
	}

	if !existingUser.MatchPassword(data.Password) {
		return utils.SendErrorResponse(
			c,
			http.StatusNoContent,
			fmt.Sprintf("The provided email '%s' may not exist or the provided password is incorrect.", data.Email),
		)
	}

	if !existingUser.Verified {
		return utils.SendErrorResponse(
			c,
			http.StatusUnauthorized,
			"You must verify your email before signing in! Check your inbox for account verification instructions "+
				"or generate another account verification email.",
		)
	}

	token, exp, err := existingUser.GenerateSessionToken()
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	cookie := fiber.Cookie{
		Name:     "SESSION_TOKEN",
		Value:    token,
		Expires:  exp,
		Path:     "/",
		HTTPOnly: true,
		//Secure: true,
	}

	c.Cookie(&cookie)
	return c.Status(http.StatusOK).Send(nil)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "SESSION_TOKEN",
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HTTPOnly: true,
		//Secure: true,
	}

	c.Cookie(&cookie)
	return c.Status(http.StatusOK).Send(nil)
}

func VerifyAccount(c *fiber.Ctx) error {
	var db = database.GetConnection()
	token := c.Query("token")
	if len(token) == 0 {
		return utils.SendErrorResponse(
			c,
			http.StatusUnauthorized,
			"You must provide a valid account verification token!",
		)
	}

	parsedToken, err := utils.ValidateUserToken(token)
	if err != nil {
		return utils.SendErrorResponse(
			c,
			http.StatusUnauthorized,
			"The provided token is not valid. If the account verification token was sent over 30 days ago, you will "+
				"need to generate another account verification email.",
		)
	}

	var user models.User
	if err := db.Where("email=?", &parsedToken.Email).First(&user).Error; err != nil {
		return utils.SendErrorResponse(c, http.StatusUnprocessableEntity, "")
	}

	if user.Verified {
		return c.Status(http.StatusNotModified).Send(nil)
	}

	var newToken []byte
	err = db.Model(&user).Updates(models.User{Verified: true, Token: &newToken}).Error
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusOK).SendString(fmt.Sprintf("Successfully verified %s!", user.Email))
}

func ResendAccountVerification(c *fiber.Ctx) error {
	var db = database.GetConnection()
	data := new(ReqUser)
	if err := c.BodyParser(data); err != nil {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "You must provide a valid email!")
	}

	// TODO(carlotta): Add field validations for "email"
	if len(data.Email) == 0 {
		return utils.SendErrorResponse(
			c,
			http.StatusBadRequest,
			"You must provide a valid email address to resend an account verification email to.",
		)
	}

	var user models.User
	if err := db.Where("email=?", &data.Email).First(&user).Error; err != nil {
		return utils.SendErrorResponse(c, http.StatusUnprocessableEntity, "")
	}

	if user.Verified {
		return c.Status(http.StatusNotModified).Send(nil)
	}

	token, _, err := utils.GenerateUserToken(data.Email)
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	err = db.Model(&user).Update("token", &token).Error
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	// TODO(carlotta): Send account verification email

	return c.Status(http.StatusAccepted).SendString(fmt.Sprintf("Resent a verification email to %s.", user.Email))
}

func SendResetPasswordEmail(c *fiber.Ctx) error {
	var db = database.GetConnection()
	data := new(ReqUser)
	if err := c.BodyParser(data); err != nil {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "You must provide a valid email!")
	}

	// TODO(carlotta): Add field validations for "email"
	if len(data.Email) == 0 {
		return utils.SendErrorResponse(
			c,
			http.StatusBadRequest,
			"You must provide a valid email address to send an password email to.",
		)
	}

	var user models.User
	if err := db.Where("email=?", &data.Email).First(&user).Error; err != nil {
		return utils.SendErrorResponse(
			c,
			http.StatusUnprocessableEntity,
			"",
		)
	}

	token, _, err := utils.GenerateUserToken(data.Email)
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	err = db.Model(&user).Update("token", &token).Error
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	// TODO(carlotta): Send account verification email

	return c.Status(http.StatusAccepted).SendString(fmt.Sprintf("Sent a password reset email to %s.", user.Email))
}

func UpdatePassword(c *fiber.Ctx) error {
	var db = database.GetConnection()
	data := new(ReqUser)
	if err := c.BodyParser(data); err != nil {
		return utils.SendErrorResponse(
			c,
			http.StatusBadRequest,
			"You must provide a valid email and password reset token!",
		)
	}

	// TODO(carlotta): Add field validations for "password" and "token"
	if len(data.Password) == 0 || len(data.Token) == 0 {
		return utils.SendErrorResponse(
			c,
			http.StatusBadRequest,
			"You must provide a valid password and password reset token.",
		)
	}

	parsedToken, err := utils.ValidateUserToken(data.Token)
	if err != nil {
		return utils.SendErrorResponse(
			c,
			http.StatusUnauthorized,
			"The provided token is not valid. If the token was sent over 30 days ago, you will need to generate "+
				"another reset password email.",
		)
	}

	var user models.User
	if err := db.Where("email=? AND token IS NOT NULL", &parsedToken.Email).First(&user).Error; err != nil {
		return utils.SendErrorResponse(
			c,
			http.StatusUnprocessableEntity,
			"",
		)
	}

	var newToken []byte
	newPassword, err := utils.CreateEncryptedPassword([]byte(data.Password))
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	err = db.Model(&user).Updates(models.User{Password: newPassword, Token: &newToken}).Error
	if err != nil {
		return utils.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.Status(http.StatusCreated).SendString("Your account has been updated with a new password!")
}

func GetAccountInfo(c *fiber.Ctx) error {
	var db = database.GetConnection()
	var userSessionId = c.Locals("userSessionId").(string)

	var user models.User
	if err := db.Where("id=?", &userSessionId).First(&user).Error; err != nil {
		return utils.SendErrorResponse(
			c,
			http.StatusInternalServerError,
			"Encountered an unexpected error. Unable to locate the associated account.",
		)
	}

	return c.Status(http.StatusOK).JSON(user)
}

func DeleteAccount(c *fiber.Ctx) error {
	var db = database.GetConnection()
	var userSessionId = c.Locals("userSessionId").(string)

	var user models.User
	if err := db.Where("id=?", &userSessionId).First(&user).Error; err != nil {
		return utils.SendErrorResponse(
			c,
			http.StatusInternalServerError,
			"Encountered an unexpected error. Unable to locate the associated account.",
		)
	}

	var err = db.Delete(&user).Error
	if err != nil {
		return utils.SendErrorResponse(
			c,
			http.StatusInternalServerError,
			"Encountered an unexpected error. Unable to delete account.",
		)
	}

	return Logout(c)
}

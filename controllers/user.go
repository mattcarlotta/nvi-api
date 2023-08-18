package controllers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	db := database.GetConnection()

	data := new(ReqUser)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid name, email and password!"},
		)
	}

	// TODO(carlotta): Add field validations for "name," "email," and "password"
	if len(data.Email) == 0 || len(data.Password) == 0 || len(data.Name) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid name, email and password!"},
		)
	}

	var user models.User
	if err := db.Where(&models.User{Email: data.Email}).First(&user).Error; err == nil {
		return c.Status(fiber.StatusOK).JSON(
			fiber.Map{"error": fmt.Sprintf(
				"The provided email '%s' may already exist or is not a valid email address.", data.Email),
			},
		)
	}

	token, _, err := utils.GenerateUserToken(data.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	newToken := []byte(token)
	newUser := models.User{
		Email:    data.Email,
		Name:     data.Name,
		Password: []byte(data.Password),
		Token:    &newToken,
	}

	if err = db.Create(&newUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// TODO(carlotta): Send account verification email

	return c.Status(fiber.StatusCreated).SendString(
		fmt.Sprintf("Welcome, %s! Please check your %s inbox for steps to verify your account.", data.Name, data.Email),
	)
}

func Login(c *fiber.Ctx) error {
	db := database.GetConnection()

	data := new(ReqUser)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid email and password!"},
		)
	}

	// TODO(carlotta): Add field validations for "email" and "password"
	if len(data.Email) == 0 || len(data.Password) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid email and password!"},
		)
	}

	var existingUser models.User
	if err := db.Where(&models.User{Email: data.Email}).First(&existingUser).Error; err != nil {
		return c.Status(fiber.StatusNoContent).JSON(
			fiber.Map{"error": fmt.Sprintf(
				"The provided email '%s' may not exist or the provided password is incorrect.", data.Email),
			},
		)
	}

	if !existingUser.MatchPassword(data.Password) {
		return c.Status(fiber.StatusNoContent).JSON(
			fiber.Map{"error": fmt.Sprintf(
				"The provided email '%s' may not exist or the provided password is incorrect.", data.Email),
			},
		)
	}

	if !existingUser.Verified {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{"error": "You must verify your email before signing in! Check your inbox for account " +
				"verification instructions or generate another account verification email.",
			},
		)
	}

	token, exp, err := existingUser.GenerateSessionToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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

	data := new(ReqUser)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid email!"},
		)
	}

	// TODO(carlotta): Add field validations for "email"
	if len(data.Email) == 0 {
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

	data := new(ReqUser)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid email!"},
		)
	}

	// TODO(carlotta): Add field validations for "email"
	if len(data.Email) == 0 {
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

	data := new(ReqUser)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "You must provide a valid password and password reset token."},
		)
	}

	// TODO(carlotta): Add field validations for "password" and "token"
	if len(data.Password) == 0 {
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

	newPassword, err := utils.CreateEncryptedPassword([]byte(data.Password))
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
	userSessionId := c.Locals("userSessionId").(uuid.UUID)

	var user models.User
	if err := db.Where(&models.User{ID: userSessionId}).First(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "Encountered an unexpected error. Unable to locate the associated account."},
		)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func DeleteAccount(c *fiber.Ctx) error {
	db := database.GetConnection()
	userSessionId := c.Locals("userSessionId").(uuid.UUID)

	var user models.User
	if err := db.Where(&models.User{ID: userSessionId}).First(&user).Error; err != nil {
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

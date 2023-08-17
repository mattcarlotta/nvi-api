package controllers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
	"github.com/mattcarlotta/nvi-api/utils"
)

type ReqUser struct {
	Name     string `json:"name" xml:"name" form:"name"`
	Email    string `json:"email" xml:"email" form:"email"`
	Password string `json:"password" xml:"password" form:"password"`
	Token    string `json:"token"`
}

func Register(c *fiber.Ctx) error {
	var db = database.GetConnection()

	// TODO(carlotta): Add field validations for "name," "email," and "password"
	data := new(ReqUser)
	if err := c.BodyParser(data); err != nil {
		return utils.SendErrorResponse(
			c,
			http.StatusBadRequest,
			"You must provide a valid email and password!",
		)
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

	return c.Status(http.StatusCreated).SendString(fmt.Sprintf("Welcome, %s! Please check your %s inbox for steps to verify your account.", data.Name, data.Email))
}

func Login(c *fiber.Ctx) error {
	var db = database.GetConnection()

	// TODO(carlotta): Add field validations for "email" and "password"
	data := new(ReqUser)
	if err := c.BodyParser(data); err != nil {
		return utils.SendErrorResponse(c, http.StatusBadRequest, "You must provide a valid name, email and password!")
	}

	var existingUser models.User
	if err := db.Where("email=?", &data.Email).First(&existingUser).Error; err != nil {
		return utils.SendErrorResponse(
			c,
			http.StatusOK,
			fmt.Sprintf("The provided email '%s' may not exist or the provided password is incorrect.", data.Email),
		)
	}

	if !existingUser.MatchPassword(data.Password) {
		return utils.SendErrorResponse(
			c,
			http.StatusOK,
			fmt.Sprintf("The provided email '%s' may not exist or the provided password is incorrect.", data.Email),
		)
	}

	if !existingUser.Verified {
		return utils.SendErrorResponse(
			c,
			http.StatusUnauthorized,
			"You must verify your email before signing in! Check your inbox for account verification instructions or generate another account verification email.",
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
	return c.Status(http.StatusCreated).Send(nil)
}

func Logout(c *fiber.Ctx) error {
	c.ClearCookie("SESSION_TOKEN")
	return c.Status(http.StatusOK).Send(nil)
}

// func VerifyAccount(res http.ResponseWriter, req *http.Request) {
// 	var db = database.GetConnection()
// 	query := req.URL.Query()
// 	token := query.Get("token")
// 	if len(token) == 0 {
// 		utils.SendErrorResponse(
// 			res,
// 			http.StatusUnauthorized,
// 			"You must provide a valid account verification token!",
// 		)
// 		return

// 	}
// 	parsedToken, err := utils.ValidateUserToken(token)
// 	if err != nil {
// 		utils.SendErrorResponse(
// 			res,
// 			http.StatusUnauthorized,
// 			"The provided token is not valid. If the token was sent over 30 days ago, you will need to generate another account verification email.",
// 		)
// 		return
// 	}

// 	var user models.User
// 	if err := db.Where("email=? AND token IS NOT NULL", &parsedToken.Email).First(&user).Error; err != nil {
// 		utils.SendErrorResponse(
// 			res,
// 			http.StatusUnprocessableEntity,
// 			"",
// 		)
// 		return
// 	}

// 	if user.Verified {
// 		res.WriteHeader(http.StatusNotModified)
// 		return
// 	}

// 	var newToken []byte
// 	err = db.Model(&user).Updates(models.User{Verified: true, Token: &newToken}).Error
// 	if err != nil {
// 		utils.SendErrorResponse(res, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	res.Header().Set("Content-Type", "text/plain")
// 	res.WriteHeader(http.StatusOK)
// 	res.Write([]byte(fmt.Sprintf("Successfully verified %s!", user.Email)))
// }

// func ResendAccountVerification(res http.ResponseWriter, req *http.Request) {
// 	var db = database.GetConnection()
// 	body, err := io.ReadAll(req.Body)
// 	if err != nil || len(body) == 0 {
// 		utils.SendErrorResponse(res, http.StatusBadRequest, "You must provide a valid email!")
// 		return
// 	}
// 	// TODO(carlotta): Add field validations for "email" and "password"
// 	var unauthedUser ReqUser
// 	err = json.Unmarshal(body, &unauthedUser)
// 	if err != nil {
// 		utils.SendErrorResponse(res, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	if len(unauthedUser.Email) == 0 {
// 		utils.SendErrorResponse(
// 			res,
// 			http.StatusBadRequest,
// 			"You must provide a valid email address to resend an account verification email to.",
// 		)
// 		return
// 	}

// 	var user models.User
// 	if err := db.Where("email=?", &unauthedUser.Email).First(&user).Error; err != nil {
// 		utils.SendErrorResponse(
// 			res,
// 			http.StatusUnprocessableEntity,
// 			"",
// 		)
// 		return
// 	}

// 	if user.Verified {
// 		res.WriteHeader(http.StatusNotModified)
// 		return
// 	}

// 	token, _, err := utils.GenerateUserToken(unauthedUser.Email)
// 	if err != nil {
// 		utils.SendErrorResponse(res, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	err = db.Model(&user).Update("token", &token).Error
// 	if err != nil {
// 		utils.SendErrorResponse(res, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	// TODO(carlotta): Send account verification email

// 	res.WriteHeader(http.StatusAccepted)
// 	res.Write([]byte(fmt.Sprintf("Resent a verification email to %s.", user.Email)))
// }

// func SendResetPasswordEmail(res http.ResponseWriter, req *http.Request) {
// 	var db = database.GetConnection()
// 	body, err := io.ReadAll(req.Body)
// 	if err != nil || len(body) == 0 {
// 		utils.SendErrorResponse(res, http.StatusBadRequest, "You must provide a valid email!")
// 		return
// 	}
// 	// TODO(carlotta): Add field validations for "email"
// 	var unauthedUser ReqUser
// 	err = json.Unmarshal(body, &unauthedUser)
// 	if err != nil {
// 		utils.SendErrorResponse(res, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	if len(unauthedUser.Email) == 0 {
// 		utils.SendErrorResponse(
// 			res,
// 			http.StatusBadRequest,
// 			"You must provide a valid email address to send an password email to.",
// 		)
// 		return
// 	}

// 	var user models.User
// 	if err := db.Where("email=?", &unauthedUser.Email).First(&user).Error; err != nil {
// 		utils.SendErrorResponse(
// 			res,
// 			http.StatusUnprocessableEntity,
// 			"",
// 		)
// 		return
// 	}

// 	token, _, err := utils.GenerateUserToken(unauthedUser.Email)
// 	if err != nil {
// 		utils.SendErrorResponse(res, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	err = db.Model(&user).Update("token", &token).Error
// 	if err != nil {
// 		utils.SendErrorResponse(res, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	// TODO(carlotta): Send account verification email

// 	res.WriteHeader(http.StatusAccepted)
// 	res.Write([]byte(fmt.Sprintf("Sent a password reset email to %s.", user.Email)))
// }

// func UpdatePassword(res http.ResponseWriter, req *http.Request) {
// 	var db = database.GetConnection()
// 	body, err := io.ReadAll(req.Body)
// 	if err != nil || len(body) == 0 {
// 		utils.SendErrorResponse(res, http.StatusBadRequest, "You must provide a valid email!")
// 		return
// 	}
// 	// TODO(carlotta): Add field validations for "password" and "token"
// 	var data ReqUser
// 	err = json.Unmarshal(body, &data)
// 	if err != nil {
// 		utils.SendErrorResponse(res, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	if len(data.Password) == 0 || len(data.Token) == 0 {
// 		utils.SendErrorResponse(
// 			res,
// 			http.StatusBadRequest,
// 			"You must provide a valid password and password reset token.",
// 		)
// 		return
// 	}

// 	parsedToken, err := utils.ValidateUserToken(data.Token)
// 	if err != nil {
// 		utils.SendErrorResponse(
// 			res,
// 			http.StatusUnauthorized,
// 			"The provided token is not valid. If the token was sent over 30 days ago, you will need to generate another reset password email.",
// 		)
// 		return
// 	}

// 	var user models.User
// 	if err := db.Where("email=? AND token IS NOT NULL", &parsedToken.Email).First(&user).Error; err != nil {
// 		utils.SendErrorResponse(
// 			res,
// 			http.StatusUnprocessableEntity,
// 			"",
// 		)
// 		return
// 	}

// 	var newToken []byte
// 	newPassword, err := utils.CreateEncryptedPassword([]byte(data.Password))
// 	if err != nil {
// 		utils.SendErrorResponse(res, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	err = db.Model(&user).Updates(models.User{Password: newPassword, Token: &newToken}).Error
// 	if err != nil {
// 		utils.SendErrorResponse(res, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	res.WriteHeader(http.StatusCreated)
// 	res.Write([]byte("Your account has been updated with a new password!"))
// }

// func DeleteAccount(res http.ResponseWriter, req *http.Request) {
// 	var db = database.GetConnection()
// 	var userSessionId = utils.GetUserSessionId(res, req)

// 	var user models.User
// 	if err := db.Where("id=?", &userSessionId).First(&user).Error; err != nil {
// 		utils.SendErrorResponse(
// 			res,
// 			http.StatusInternalServerError,
// 			"Encountered an unexpected error. Unable to locate the associated account.",
// 		)
// 		return
// 	}

// 	var err = db.Delete(&user).Error
// 	if err != nil {
// 		utils.SendErrorResponse(
// 			res,
// 			http.StatusInternalServerError,
// 			"Encountered an unexpected error. Unable to delete account.",
// 		)
// 		return
// 	}

// 	Logout(res, req)
// }

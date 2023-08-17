package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
	"github.com/mattcarlotta/nvi-api/utils"
)

type ReqUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(res http.ResponseWriter, req *http.Request) {
	var db = database.GetConnection()

	var data ReqUser
	body, err := io.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		utils.SendErrorResponse(res, http.StatusBadRequest, "You must provide a valid name, email and password!")
		return
	}

	// TODO(carlotta): Add field validations for "name," "email," and "password"
	err = json.Unmarshal(body, &data)
	if err != nil {
		utils.SendErrorResponse(res, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	if err := db.Where("email=?", &data.Email).First(&user).Error; err == nil {
		utils.SendErrorResponse(
			res,
			http.StatusOK,
			fmt.Sprintf("The provided email '%s' may already exist or is not a valid email address.", data.Email),
		)
		return
	}

	token, _, err := utils.GenerateUserToken(data.Email)
	if err != nil {
		utils.SendErrorResponse(res, http.StatusInternalServerError, err.Error())
		return
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
		utils.SendErrorResponse(res, http.StatusInternalServerError, err.Error())
		return
	}

	// TODO(carlotta): Send account verification email

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(fmt.Sprintf("Welcome, %s! Please check your %s inbox for steps to verify your account.", data.Name, data.Email)))
}

func Login(res http.ResponseWriter, req *http.Request) {
	var db = database.GetConnection()
	body, err := io.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		utils.SendErrorResponse(res, http.StatusBadRequest, "You must provide a valid email and password!")
		return
	}

	// TODO(carlotta): Add field validations for "email" and "password"
	var unauthedUser ReqUser
	err = json.Unmarshal(body, &unauthedUser)
	if err != nil {
		utils.SendErrorResponse(res, http.StatusBadRequest, err.Error())
		return
	}

	var existingUser models.User
	if err = db.Where("email=?", &unauthedUser.Email).First(&existingUser).Error; err != nil {
		utils.SendErrorResponse(
			res,
			http.StatusOK,
			fmt.Sprintf("The provided email '%s' may not exist or the provided password is incorrect.", unauthedUser.Email),
		)
		return
	}

	if !existingUser.MatchPassword(unauthedUser.Password) {
		utils.SendErrorResponse(
			res,
			http.StatusOK,
			fmt.Sprintf("The provided email '%s' may not exist or the provided password is incorrect.", unauthedUser.Email),
		)
		return
	}

	if !existingUser.Verified {
		utils.SendErrorResponse(
			res,
			http.StatusUnauthorized,
			"You must verify your email before signing in! Check your inbox for account verification instructions or generate another account verification email.",
		)
		return

	}

	token, exp, err := existingUser.GenerateSessionToken()
	if err != nil {
		utils.SendErrorResponse(res, http.StatusInternalServerError, err.Error())
		return
	}

	cookie := http.Cookie{
		Name:     "SESSION_TOKEN",
		Value:    token,
		Expires:  exp,
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(res, &cookie)
	res.WriteHeader(http.StatusCreated)
}

func Logout(res http.ResponseWriter, req *http.Request) {
	cookie := http.Cookie{
		Name:     "SESSION_TOKEN",
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(res, &cookie)
	res.WriteHeader(http.StatusOK)
}

func DeleteAccount(res http.ResponseWriter, req *http.Request) {
	var db = database.GetConnection()
	var userSessionId = utils.GetUserSessionId(res, req)

	var user models.User
	if err := db.Where("id=?", &userSessionId).First(&user).Error; err != nil {
		utils.SendErrorResponse(
			res,
			http.StatusInternalServerError,
			"Encountered an unexpected error. Unable to locate the associated account.",
		)
		return
	}

	var err = db.Delete(&user).Error
	if err != nil {
		utils.SendErrorResponse(
			res,
			http.StatusInternalServerError,
			"Encountered an unexpected error. Unable to delete account.",
		)
		return
	}

	Logout(res, req)
}

func VerifyAccount(res http.ResponseWriter, req *http.Request) {
	var db = database.GetConnection()
	query := req.URL.Query()
	token := query.Get("token")
	if len(token) == 0 {
		utils.SendErrorResponse(
			res,
			http.StatusUnauthorized,
			"You must provide a valid account verification token!",
		)
		return

	}
	parsedToken, err := utils.ValidateUserToken(token)
	if err != nil {
		utils.SendErrorResponse(
			res,
			http.StatusUnauthorized,
			"The provided token is not valid. If the token was sent over 30 days ago, you will need to generate another account verification email",
		)
		return
	}

	var user models.User
	if err := db.Where("email=?", &parsedToken.Email).First(&user).Error; err != nil {
		utils.SendErrorResponse(
			res,
			http.StatusUnprocessableEntity,
			"",
		)
		return
	}

	if user.Verified {
		res.WriteHeader(http.StatusNotModified)
		return
	}

	var newToken []byte
	err = db.Model(&user).Updates(models.User{Verified: true, Token: &newToken}).Error
	if err != nil {
		utils.SendErrorResponse(res, http.StatusInternalServerError, err.Error())
		return
	}

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(fmt.Sprintf("Successfully verified %s!", user.Email)))
}

func ResendAccountVerificatin(res http.ResponseWriter, req *http.Request) {
	var db = database.GetConnection()
	body, err := io.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		utils.SendErrorResponse(res, http.StatusBadRequest, "You must provide a valid email and password!")
		return
	}
	// TODO(carlotta): Add field validations for "email" and "password"
	var unauthedUser ReqUser
	err = json.Unmarshal(body, &unauthedUser)
	if err != nil {
		utils.SendErrorResponse(res, http.StatusBadRequest, err.Error())
		return
	}

	if len(unauthedUser.Email) == 0 {
		utils.SendErrorResponse(
			res,
			http.StatusBadRequest,
			"You must provide a valid email address to resend an account verification email to.",
		)
		return
	}

	var user models.User
	if err := db.Where("email=?", &unauthedUser.Email).First(&user).Error; err != nil {
		utils.SendErrorResponse(
			res,
			http.StatusUnprocessableEntity,
			"",
		)
		return
	}

	if user.Verified {
		res.WriteHeader(http.StatusNotModified)
		return
	}

	token, _, err := utils.GenerateUserToken(unauthedUser.Email)
	if err != nil {
		utils.SendErrorResponse(res, http.StatusInternalServerError, err.Error())
		return
	}

	err = db.Model(&user).Update("token", &token).Error
	if err != nil {
		utils.SendErrorResponse(res, http.StatusInternalServerError, err.Error())
		return
	}

	// TODO(carlotta): Send account verification email

	res.WriteHeader(http.StatusCreated)
}

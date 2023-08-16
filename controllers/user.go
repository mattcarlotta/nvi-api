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
			fmt.Sprintf("The provided email '%s' may already exist or is not using a valid email domain!", data.Email),
		)
		return
	}

	newUser := models.User{Email: data.Email, Name: data.Name, Password: []byte(data.Password)}
	err = db.Model(&user).Create(&newUser).Error
	if err != nil {
		utils.SendErrorResponse(res, http.StatusInternalServerError, err.Error())
		return
	}

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(fmt.Sprintf("Successfully registered %s. Welcome, %s!", data.Email, data.Name)))
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
			fmt.Sprintf("The provided email '%s' may not exist or the provided password is incorrect!", unauthedUser.Email),
		)
		return
	}

	if !existingUser.MatchPassword(unauthedUser.Password) {
		utils.SendErrorResponse(
			res,
			http.StatusOK,
			fmt.Sprintf("The provided email '%s' may not exist or the provided lassword is incorrect!", unauthedUser.Email),
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

	res.WriteHeader(http.StatusAccepted)
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

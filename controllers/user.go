package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
	"github.com/mattcarlotta/nvi-api/utils"
)

type ReqUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUser(res http.ResponseWriter, req *http.Request) {
	var db = database.GetDB()

	var newUser ReqUser
	body, err := io.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		utils.SendErrorResponse(res, http.StatusBadRequest, "You must provide a valid name, email and password!")
		return
	}

	err = json.Unmarshal(body, &newUser)
	if err != nil {
		utils.SendErrorResponse(res, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	if err := db.Where("email=?", &newUser.Email).First(&user).Error; err == nil {
		utils.SendErrorResponse(
			res,
			http.StatusOK,
			fmt.Sprintf("The provided email '%s' may already exist or is not using a valid email domain!", newUser.Email),
		)
		return
	}

	user.Email = newUser.Email
	user.Name = newUser.Name
	user.Password = []byte(newUser.Password)

	err = db.Create(&user).Error
	if err != nil {
		utils.SendErrorResponse(res, http.StatusInternalServerError, err.Error())
		return
	}

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(fmt.Sprintf("Successfully registered %s. Welcome, %s!", user.Email, user.Name)))
}

func Login(res http.ResponseWriter, req *http.Request) {
	var db = database.GetDB()
	body, err := io.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		utils.SendErrorResponse(res, http.StatusBadRequest, "You must provide a valid email and password!")
		return
	}

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
			fmt.Sprintf("The provided email '%s' may not exist or the provided password is incorrect!", unauthedUser.Email),
		)
		return
	}

	exp, token, err := existingUser.GenerateSessionToken()
	if err != nil {
		utils.SendErrorResponse(res, http.StatusInternalServerError, err.Error())
		return
	}

	cookie := http.Cookie{
		Name:    "SESSION_TOKEN",
		Value:   token,
		Expires: exp,
		Path:    "/",
	}

	http.SetCookie(res, &cookie)

	res.WriteHeader(http.StatusAccepted)
	res.Write(nil)
}

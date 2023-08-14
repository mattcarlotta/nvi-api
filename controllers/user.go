package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
	"github.com/mattcarlotta/nvi-api/utils"
)

func AllUsers(res http.ResponseWriter, req *http.Request) {
	fmt.Fprint(res, "All users endpoint\n")
}

func CreateUser(res http.ResponseWriter, req *http.Request) {
	if req.Body == nil {
		utils.SendErrorResponse(res, http.StatusBadRequest, "You must send a valid name, email and password!")
		return
	}

	var newUser models.NewUser
	err := json.NewDecoder(req.Body).Decode(&newUser)
	if err != nil {
		utils.SendErrorResponse(res, http.StatusBadRequest, fmt.Sprint(err))
		return
	}

	var user models.User
	if err := database.DB.Where("email= ?", &newUser.Email).First(&user).Error; err == nil {
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

	err = database.DB.Create(&user).Error
	if err != nil {
		utils.SendErrorResponse(res, http.StatusInternalServerError, fmt.Sprint(err))
		return
	}

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	fmt.Fprintf(res, "Successfully registered %s. Welcome, %s!", user.Email, user.Name)
}

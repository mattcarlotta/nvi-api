package controllers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
	"github.com/mattcarlotta/nvi-api/utils"
)

func CreateEnvironment(res http.ResponseWriter, req *http.Request) {
	var db = database.GetDB()
	params := mux.Vars(req)
	envName := params["name"]

	if len(envName) == 0 {
		utils.SendErrorResponse(res, http.StatusBadRequest, "You must provide a valid enviroment!")
		return
	}

	var userSessionId = utils.GetUserSessionId(res, req)

	var environment models.Environment
	if err := db.Where("name=? AND user_id=?", &envName, &userSessionId).First(&environment).Error; err == nil {
		utils.SendErrorResponse(
			res,
			http.StatusOK,
			fmt.Sprintf("The provided environment '%s' already exists. Please choose a different environment name!", envName),
		)
		return
	}

	environment.Name = envName
	environment.UserId = uuid.MustParse(userSessionId)

	var err = db.Create(&environment).Error
	if err != nil {
		utils.SendErrorResponse(res, http.StatusInternalServerError, err.Error())
		return
	}

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(fmt.Sprintf("Successfully created the %s environment!", envName)))
}

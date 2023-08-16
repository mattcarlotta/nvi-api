package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
	"github.com/mattcarlotta/nvi-api/utils"
)

type ReqEnv struct {
	OriginalName string `json:"originalName"`
	UpdatedName  string `json:"updatedName"`
}

func CreateEnvironment(res http.ResponseWriter, req *http.Request) {
	var db = database.GetDB()
	params := mux.Vars(req)
	envName := params["name"]

	if len(envName) == 0 {
		utils.SendErrorResponse(res, http.StatusBadRequest, "You must provide a valid environment!")
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

func DeleteEnvironment(res http.ResponseWriter, req *http.Request) {
	var db = database.GetDB()
	params := mux.Vars(req)
	envName := params["name"]

	if len(envName) == 0 {
		utils.SendErrorResponse(res, http.StatusBadRequest, "You must provide a valid environment!")
		return
	}

	var userSessionId = utils.GetUserSessionId(res, req)

	var environment models.Environment
	if err := db.Where("name=? AND user_id=?", &envName, &userSessionId).First(&environment).Error; err != nil {
		utils.SendErrorResponse(
			res,
			http.StatusOK,
			fmt.Sprintf("The provided environment '%s' doesn't appear to exist!", envName),
		)
		return
	}

	var err = db.Delete(&environment).Error
	if err != nil {
		utils.SendErrorResponse(res, http.StatusInternalServerError, err.Error())
		return
	}

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(fmt.Sprintf("Successfully deleted the %s environment!", envName)))
}

func UpdateEnvironment(res http.ResponseWriter, req *http.Request) {
	var db = database.GetDB()

	body, err := io.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		utils.SendErrorResponse(
			res,
			http.StatusBadRequest,
			"You must provide a valid original environment name and an updated environment name!",
		)
		return
	}

	var newEnvironment ReqEnv
	// TODO(carlotta): Add field validations for "originalName" and "updatedName"
	err = json.Unmarshal(body, &newEnvironment)
	if err != nil {
		utils.SendErrorResponse(res, http.StatusBadRequest, err.Error())
		return
	}

	var userSessionId = utils.GetUserSessionId(res, req)

	var environment models.Environment
	if err := db.Where("name=? AND user_id=?", &newEnvironment.OriginalName, &userSessionId).First(&environment).Error; err != nil {
		utils.SendErrorResponse(
			res,
			http.StatusOK,
			fmt.Sprintf("The provided environment '%s' doesn't appear to exist!", newEnvironment.OriginalName),
		)
		return
	}

	environment.Name = newEnvironment.UpdatedName
	db.Save(&environment)

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(fmt.Sprintf("Successfully updated the environment from '%s' to '%s'!", newEnvironment.OriginalName, newEnvironment.UpdatedName)))
}

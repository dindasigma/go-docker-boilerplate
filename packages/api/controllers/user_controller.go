package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/dindasigma/go-docker-boilerplate/packages/api/auth"
	"github.com/dindasigma/go-docker-boilerplate/packages/api/datasources"
	"github.com/dindasigma/go-docker-boilerplate/packages/api/models/users"
	"github.com/dindasigma/go-docker-boilerplate/packages/api/responses"
	"github.com/dindasigma/go-docker-boilerplate/packages/api/utils/formaterror"
	"github.com/gorilla/mux"
)

var (
	UserController userControllerInterface = &userController{}
)

type userControllerInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type userController struct{}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the input payload
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body users.User true "Create user"
// @Success 200 {object} users.User
// @Failure 422 {object} responses.Error
// @Failure 500 {object} responses.Error
// @Router /users [post]
func (c *userController) Create(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	user := users.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("")
	if err != nil {
		// 422
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	userCreated, err := user.Save(datasources.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		// 500
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	responses.JSON(w, http.StatusCreated, userCreated)
}

// GetUsers godoc
// @Summary Get details of all users
// @Description Get details of all users
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} users.User
// @Failure 500 {object} responses.Error
// @Router /users [get]
func (c *userController) Get(w http.ResponseWriter, r *http.Request) {
	user := users.User{}
	users, err := user.FindAll(datasources.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, users)
}

// GetUser godoc
// @Summary Get single row data of user with given id
// @Description Get single row data of user with given id
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "ID of the user"
// @Success 200 {object} users.User
// @Failure 500 {object} responses.Error
// @Router /users/{id} [get]
func (c *userController) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	user := users.User{}
	userGotten, err := user.FindByID(datasources.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusOK, userGotten)
}

// UpdateUser godoc
// @Summary Update user identified by the given id
// @Description Update user corresponding to the input id
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "ID of the user to be updated"
// @Param user body users.User true "Update user"
// @Security ApiKeyAuth
// @Success 200 {object} users.User
// @Failure 400 {object} responses.Error
// @Failure 401 {object} responses.Error
// @Failure 422 {object} responses.Error
// @Failure 500 {object} responses.Error
// @Router /users/{id} [put]
func (c *userController) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		// 400
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := users.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		// 401
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if tokenID != uint32(uid) {
		// 401
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		// 422
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	updatedUser, err := user.Update(datasources.DB, uint32(uid))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, updatedUser)
}

// DeleteUser godoc
// @Summary Delete user identified by the given id
// @Description Delete user corresponding to the input id
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "ID of the user to be deleted"
// @Security ApiKeyAuth
// @Success 204 "No Content"
// @Failure 400 {object} responses.Error
// @Failure 401 {object} responses.Error
// @Failure 500 {object} responses.Error
// @Router /users/{id} [delete]
func (c *userController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := users.User{}

	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if tokenID != 0 && tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	_, err = user.Delete(datasources.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	responses.JSON(w, http.StatusNoContent, "")
}

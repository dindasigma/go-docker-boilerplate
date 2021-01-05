package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/dindasigma/go-docker-boilerplate/packages/api/auth"
	"github.com/dindasigma/go-docker-boilerplate/packages/api/datasources"
	"github.com/dindasigma/go-docker-boilerplate/packages/api/models/users"
	"github.com/dindasigma/go-docker-boilerplate/packages/api/responses"
	"github.com/dindasigma/go-docker-boilerplate/packages/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
)

var (
	LoginController loginControllerInterface = &loginController{}
)

type loginControllerInterface interface {
	Login(w http.ResponseWriter, r *http.Request)
	SignIn(email, password string) (string, error)
}

type loginController struct{}

// @Summary Login user to get JWT Token
// @Description Login user to get JWT Token for bearerAuth
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body users.User true "Email and Password only"
// @Success 200 {string} Token "JWT Token"
// @Failure 422 {object} responses.Error
// @Router /login [post]
func (c *loginController) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// 422
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := users.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		// 422
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		// 422
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := c.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		// 422
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, token)
}

func (c *loginController) SignIn(email, password string) (string, error) {

	var err error

	user := users.User{}
	err = user.CheckUser(datasources.DB, email)
	if err != nil {
		return "", err
	}
	err = users.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(user.ID)
}

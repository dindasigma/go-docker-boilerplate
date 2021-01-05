package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/dindasigma/go-docker-boilerplate/packages/api/auth"
	"github.com/dindasigma/go-docker-boilerplate/packages/api/models/users"
	"github.com/dindasigma/go-docker-boilerplate/packages/api/responses"
	"github.com/dindasigma/go-docker-boilerplate/packages/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Login user to get JWT Token
// @Description Login user to get JWT Token for bearerAuth
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body users.User true "Email and Password only"
// @Success 200 {string} Token "JWT Token"
// @Failure 422 {object} responses.Error
// @Router /login [post]
func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
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
	token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		// 422
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, token)
}

func (server *Server) SignIn(email, password string) (string, error) {

	var err error

	user := users.User{}

	err = server.DB.Debug().Model(users.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = users.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(user.ID)
}

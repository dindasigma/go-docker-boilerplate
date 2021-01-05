package controllers

import (
	"net/http"

	"github.com/dindasigma/go-docker-boilerplate/packages/api/utils/responses"
)

var (
	HomeController homeControllerInterface = &homeController{}
)

type homeControllerInterface interface {
	Home(w http.ResponseWriter, r *http.Request)
}

type homeController struct{}

func (c *homeController) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome to The Machine")
}

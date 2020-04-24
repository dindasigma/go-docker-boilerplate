package controllers

import (
	"net/http"

	"github.com/dindasigma/go-docker-boilerplate/packages/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To This Awesome API")

}
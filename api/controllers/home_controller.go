package controllers

import (
	"net/http"

	"github.com/Stephen00/go-fullstack/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome to this awesome api")
}

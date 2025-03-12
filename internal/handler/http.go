package handler

import (
	"github/arshamroshannejad/task-rootext/internal/helpers"
	"net/http"
)

func HttpMethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	body := "The requested resource could not be found. Please check the URL and try again."
	helpers.WriteJson(w, http.StatusNotFound, helpers.M{"error": body})
}

func HttpRequestNotFound(w http.ResponseWriter, r *http.Request) {
	body := "The requested resource could not be found. Please check the URL and try again."
	helpers.WriteJson(w, http.StatusNotFound, helpers.M{"error": body})
}

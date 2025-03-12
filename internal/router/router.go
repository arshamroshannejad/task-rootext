package router

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func SetupRoutes() http.Handler {
	r := chi.NewRouter()
	return r
}

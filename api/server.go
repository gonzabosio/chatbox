package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type App struct {
	router *chi.Mux
}

func (a *App) InitServer() {
	a.router = chi.NewRouter()
	a.router.Post("/login", login)
}

func (a *App) Run() error {
	return http.ListenAndServe(":8000", a.router)
}

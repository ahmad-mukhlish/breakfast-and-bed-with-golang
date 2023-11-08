package main

import (
	"net/http"

	"github.com/ahmad-mukhlish/annahdloh-landing-page/pkg/config"
	"github.com/ahmad-mukhlish/annahdloh-landing-page/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func handleRoute(appConfig *config.AppConfig) http.Handler {

	router := chi.NewRouter()

	router.Use(NoSurf)
	router.Use(middleware.Recoverer)
	router.Use(SessionLoad)

	router.Get("/", handlers.Repo.Home)
	router.Get("/about", handlers.Repo.About)

	return router
}

package main

import (
	"net/http"

	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/pkg/config"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/pkg/handlers"
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
	router.Get("/major", handlers.Repo.Major)
	router.Get("/general", handlers.Repo.General)
	router.Get("/contact", handlers.Repo.Contact)
	router.Get("/reservation", handlers.Repo.Reservation)
	router.Get("/check-availability", handlers.Repo.CheckAvailability)
	router.Post("/check-availability", handlers.Repo.PostCheckAvailability)

	rootDirectoryStaticFile := http.Dir("./res/")
	staticFileServer := http.FileServer(rootDirectoryStaticFile)

	resPath := appConfig.ResRoutePath

	router.Handle(resPath+"*", http.StripPrefix(resPath, staticFileServer))

	return router
}

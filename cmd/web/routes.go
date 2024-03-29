package main

import (
	"net/http"

	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func HandleRoute() http.Handler {

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
	router.Post("/reservation", handlers.Repo.PostReservation)
	router.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	router.Get("/check-availability", handlers.Repo.CheckAvailability)
	router.Post("/check-availability", handlers.Repo.PostCheckAvailability)
	router.Post("/check-availability/json", handlers.Repo.PostCheckAvailabilityJSON)

	router.Get("/check/room/{id}", handlers.Repo.CheckRoom)

	router.Get("/user/login", handlers.Repo.Login)
	router.Post("/user/login", handlers.Repo.PostLogin)

	router.Get("/user/logout", handlers.Repo.Logout)

	rootDirectoryStaticFile := http.Dir("./res/")
	staticFileServer := http.FileServer(rootDirectoryStaticFile)

	//admin secure pages
	router.Route("/admin", func(adminRoute chi.Router) {

		adminRoute.Use(Auth)

		adminRoute.Get("/dashboard", handlers.Repo.AdminDashboard)

	})

	router.Handle("/res"+"*", http.StripPrefix("/res", staticFileServer))

	return router
}

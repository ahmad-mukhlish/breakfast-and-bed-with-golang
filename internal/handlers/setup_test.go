package handlers

import (
	"encoding/gob"
	appConfig "github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/config"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/model"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/renders"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"net/http"
	"time"
)

var AppConfig *appConfig.AppConfig

func getRoutes() http.Handler {

	err := setupServer()
	if err != nil {
		return nil
	}

	return HandleRoute()

}

func setupServer() error {
	var err error

	AppConfig, err = setupConfig()
	if err != nil {
		return err
	}

	setupSession()
	setupRepository()
	setupLogger()

	return err

}

func setupConfig() (*appConfig.AppConfig, error) {
	var app appConfig.AppConfig

	app.UseCache = false
	app.IsProductionMode = false
	renders.PathToTemplate = "./../../templates"
	templateCache, err := renders.CreateTemplateCache()
	if err != nil {
		return &app, err
	}

	renders.SetConfig(&app)
	app.TemplateCache = templateCache

	return &app, nil
}

func setupSession() {
	session := scs.New()

	session.Lifetime = 24 * time.Hour
	session.Cookie.Secure = AppConfig.IsProductionMode
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode

	//register custom types here with gob.register
	gob.Register(model.Reservation{})

	AppConfig.Session = session
}

func setupRepository() {
	//TODO @ahmad-mukhlis to be fixed
	repo := CreateRepository(AppConfig)
	CreateHandlers(repo)
}

func HandleRoute() http.Handler {

	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(SessionLoad)

	router.Get("/", Repo.Home)
	router.Get("/about", Repo.About)

	router.Get("/major", Repo.Major)
	router.Get("/general", Repo.General)
	router.Get("/contact", Repo.Contact)
	router.Get("/reservation", Repo.Reservation)
	router.Post("/reservation", Repo.PostReservation)
	router.Get("/reservation-summary", Repo.ReservationSummary)

	router.Get("/check-availability", Repo.CheckAvailability)
	router.Post("/check-availability", Repo.PostCheckAvailability)
	router.Post("/check-availability/json", Repo.CheckAvailabilityJSON)

	rootDirectoryStaticFile := http.Dir("./res/")
	staticFileServer := http.FileServer(rootDirectoryStaticFile)

	router.Handle("/res"+"*", http.StripPrefix("/res", staticFileServer))

	return router
}

func CreateCookie() http.Cookie {
	return http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
}

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	cookie := CreateCookie()

	csrfHandler.SetBaseCookie(cookie)
	return csrfHandler

}

func SessionLoad(next http.Handler) http.Handler {
	return AppConfig.Session.LoadAndSave(next)
}

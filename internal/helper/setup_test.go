package helper

import (
	"encoding/gob"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/config"

	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/model"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/renders"
	"github.com/alexedwards/scs/v2"
	"github.com/justinas/nosurf"
	"log"
	"net/http"
	"os"
	"time"
)

//func getRoutes() http.Handler {
//
//	err := setupServer()
//	if err != nil {
//		return nil
//	}
//
//	return HandleRoute()
//
//}

func setupServer() error {
	var err error

	appConfig, err = setupConfig()
	if err != nil {
		return err
	}

	setupSession()
	//setupRepository()
	setupLogger()

	return err

}

func setupConfig() (*config.AppConfig, error) {
	var app config.AppConfig

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
	session.Cookie.Secure = appConfig.IsProductionMode
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode

	//register custom types here with gob.register
	gob.Register(model.Reservation{})

	appConfig.Session = session
}

//func setupRepository() {
//	repo := handlers.CreateRepository(appConfig)
//	handlers.CreateHandlers(repo)
//}

//func HandleRoute() http.Handler {
//
//	router := chi.NewRouter()
//
//	router.Use(middleware.Recoverer)
//	router.Use(SessionLoad)
//
//	router.Get("/", handlers.Repo.Home)
//	router.Get("/about", handlers.Repo.About)
//
//	router.Get("/major", handlers.Repo.Major)
//	router.Get("/general", handlers.Repo.General)
//	router.Get("/contact", handlers.Repo.Contact)
//	router.Get("/reservation", handlers.Repo.Reservation)
//	router.Post("/reservation", handlers.Repo.PostReservation)
//	router.Get("/reservation-summary", handlers.Repo.ReservationSummary)
//
//	router.Get("/check-availability", handlers.Repo.CheckAvailabilityForRoomById)
//	router.Post("/check-availability", handlers.Repo.PostCheckAvailability)
//	router.Post("/check-availability/json", handlers.Repo.CheckAvailabilityJSON)
//
//	rootDirectoryStaticFile := http.Dir("./res/")
//	staticFileServer := http.FileServer(rootDirectoryStaticFile)
//
//	router.Handle("/res"+"*", http.StripPrefix("/res", staticFileServer))
//
//	return router
//}

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
	return appConfig.Session.LoadAndSave(next)
}

func setupLogger() {
	appConfig.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	appConfig.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

}

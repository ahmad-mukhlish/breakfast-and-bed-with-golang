package main

import (
	"encoding/gob"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/model"
	"log"
	"net/http"
	"time"

	appConfig "github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/config"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/handlers"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/renders"
	"github.com/alexedwards/scs/v2"
)

var AppConfig *appConfig.AppConfig

func main() {
	err := setupServer()
	if err != nil {
		log.Fatal(err)
	}
	_, err = startServer(":8080")

}

func setupServer() error {
	var err error

	AppConfig, err = setupConfig()
	if err != nil {
		return err
	}

	setupSession()
	setupRepository()

	return err

}

func startServer(port string) (*http.Server, error) {

	handledRoutes := HandleRoute()

	server := &http.Server{
		Addr:    port,
		Handler: handledRoutes,
	}

	err := server.ListenAndServe()

	return server, err
}

func setupConfig() (*appConfig.AppConfig, error) {
	var app appConfig.AppConfig

	app.UseCache = false
	app.IsProductionMode = false

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
	repo := handlers.CreateRepository(AppConfig)
	handlers.CreateHandlers(repo)
}

package main

import (
	"encoding/gob"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/driver"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/helper"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/model"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/config"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/handlers"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/renders"
	"github.com/alexedwards/scs/v2"
)

var AppConfig *config.AppConfig

func main() {
	db, err := setupServer()
	defer db.DbPool.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected to db")
	_, err = startServer(":8080")

}

func setupServer() (*driver.DB, error) {
	var err error

	AppConfig, err = setupConfig()
	if err != nil {
		return nil, err
	}

	setupSession()
	setupRepository()
	setupLogger()
	db, err := setupDB()
	if err != nil {
		return nil, err
	}

	return db, err

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

func setupConfig() (*config.AppConfig, error) {
	var app config.AppConfig

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

func setupLogger() {
	AppConfig.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	AppConfig.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	helper.SetHelperAppConfig(AppConfig)
}

func setupDB() (*driver.DB, error) {
	log.Println("connecting to db")
	dbPath := "host=localhost port=54321 dbname=bookings user=ahmadmukhlis password=password"
	db, err := driver.ConnectSQL(dbPath)
	if err != nil {
		log.Fatal("error", err)
		return nil, err
	}

	return db, nil

}

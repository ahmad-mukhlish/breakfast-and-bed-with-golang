package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/pkg/config"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/pkg/model"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/pkg/renders"
)

const IPAddressKey = "ip_address"

type Repository struct {
	AppConfig *config.AppConfig
}

var Repo *Repository

func CreateRepository(appConfig *config.AppConfig) *Repository {

	return &Repository{
		AppConfig: appConfig,
	}

}

func CreateHandlers(repository *Repository) {
	Repo = repository
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {

	initializedTempalte := initiateTemplate(repo.AppConfig, r.Context())
	renders.ServeTemplate(w, r, "about.page.tmpl", initializedTempalte)
}

func (repo *Repository) General(w http.ResponseWriter, r *http.Request) {

	initializedTempalte := initiateTemplate(repo.AppConfig, r.Context())
	renders.ServeTemplate(w, r, "general.page.tmpl", initializedTempalte)
}

func (repo *Repository) Major(w http.ResponseWriter, r *http.Request) {

	initializedTempalte := initiateTemplate(repo.AppConfig, r.Context())
	renders.ServeTemplate(w, r, "major.page.tmpl", initializedTempalte)
}

func (repo *Repository) Contact(w http.ResponseWriter, r *http.Request) {

	initializedTempalte := initiateTemplate(repo.AppConfig, r.Context())
	renders.ServeTemplate(w, r, "contact.page.tmpl", initializedTempalte)
}

func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {

	IPAddrress := r.RemoteAddr
	repo.AppConfig.Session.Put(r.Context(), IPAddressKey, IPAddrress)

	initializedTempalte := initiateTemplate(repo.AppConfig, r.Context())
	renders.ServeTemplate(w, r, "home.page.tmpl", initializedTempalte)

}

func (repo *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	IPAddrress := r.RemoteAddr
	repo.AppConfig.Session.Put(r.Context(), IPAddressKey, IPAddrress)

	initializedTempalte := initiateTemplate(repo.AppConfig, r.Context())
	renders.ServeTemplate(w, r, "reservation.page.tmpl", initializedTempalte)

}

func (repo *Repository) CheckAvailability(w http.ResponseWriter, r *http.Request) {

	IPAddrress := r.RemoteAddr
	repo.AppConfig.Session.Put(r.Context(), IPAddressKey, IPAddrress)

	initializedTempalte := initiateTemplate(repo.AppConfig, r.Context())
	renders.ServeTemplate(w, r, "check-availability.page.tmpl", initializedTempalte)

}

func (repo *Repository) PostCheckAvailability(w http.ResponseWriter, r *http.Request) {

	arrival := r.Form.Get("start")
	departure := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("Your arrival date is %s, your departure date is %s", arrival, departure)))
}

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func (repo *Repository) CheckAvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	response := jsonResponse{Ok: true, Message: "Hello JSON brader"}

	ouput, err := json.MarshalIndent(response, "", "  ")

	if err != nil {
		log.Println("error")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(ouput)
}

func initiateTemplate(
	appConfig *config.AppConfig,
	context context.Context) *model.TemplateData {

	stringMap := map[string]string{}

	stringMap["test"] = "this is some string"
	stringMap[IPAddressKey] = appConfig.Session.GetString(context, IPAddressKey)

	templateData := model.TemplateData{
		StringMap: stringMap,
	}

	return &templateData

}

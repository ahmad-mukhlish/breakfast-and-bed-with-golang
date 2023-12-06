package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/config"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/form"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/model"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/renders"
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

	initializedTemplate := initiateTemplate(repo.AppConfig, r.Context())
	renders.ServeTemplate(w, r, "about.page.tmpl", initializedTemplate)
}

func (repo *Repository) General(w http.ResponseWriter, r *http.Request) {

	initializedTemplate := initiateTemplate(repo.AppConfig, r.Context())
	renders.ServeTemplate(w, r, "general.page.tmpl", initializedTemplate)
}

func (repo *Repository) Major(w http.ResponseWriter, r *http.Request) {

	initializedTemplate := initiateTemplate(repo.AppConfig, r.Context())
	renders.ServeTemplate(w, r, "major.page.tmpl", initializedTemplate)
}

func (repo *Repository) Contact(w http.ResponseWriter, r *http.Request) {

	initializedTemplate := initiateTemplate(repo.AppConfig, r.Context())
	renders.ServeTemplate(w, r, "contact.page.tmpl", initializedTemplate)
}

func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {

	initializedTemplate := initiateTemplate(repo.AppConfig, r.Context())
	renders.ServeTemplate(w, r, "home.page.tmpl", initializedTemplate)

}

func (repo *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	initializedTemplate := initiateTemplate(repo.AppConfig, r.Context())
	renders.ServeTemplate(w, r, "reservation.page.tmpl", initializedTemplate)

}

func (repo *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	reservation := model.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	actualForm := form.New(r.PostForm)

	actualForm.ValidateLength("first_name", 3)
	actualForm.Required("first_name", "last_name", "phone")
	actualForm.ValidateEmail("email")

	if !actualForm.IsValid() {

		data := make(map[string]interface{})
		data["reservation"] = reservation

		renders.ServeTemplate(w, r, "reservation.page.tmpl", &model.TemplateData{
			Data:          data,
			FormValidator: actualForm,
		})

		return

	} else {
		http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
		return
	}
}

func (repo *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {

	initializedTemplate := initiateTemplate(repo.AppConfig, r.Context())
	renders.ServeTemplate(w, r, "reservation-summary.page.tmpl", initializedTemplate)

}

func (repo *Repository) CheckAvailability(w http.ResponseWriter, r *http.Request) {

	initializedTemplate := initiateTemplate(repo.AppConfig, r.Context())
	renders.ServeTemplate(w, r, "check-availability.page.tmpl", initializedTemplate)

}

func (repo *Repository) PostCheckAvailability(w http.ResponseWriter, r *http.Request) {

	arrival := r.Form.Get("start")
	departure := r.Form.Get("end")

	_, err := w.Write([]byte(fmt.Sprintf("Your arrival date is %s, your departure date is %s", arrival, departure)))
	if err != nil {
		return
	}
}

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func (repo *Repository) CheckAvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	arrival := r.Form.Get("start")
	departure := r.Form.Get("end")
	message := fmt.Sprintf("Your arrival date is %s, your departure date is %s", arrival, departure)

	response := jsonResponse{Ok: true, Message: message}
	output, err := json.MarshalIndent(response, "", "  ")

	if err != nil {
		log.Println("error")
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		return
	}
}

func initiateTemplate(
	appConfig *config.AppConfig,
	context context.Context) *model.TemplateData {

	stringMap := map[string]string{}

	stringMap["test"] = "this is some string"
	stringMap[IPAddressKey] = appConfig.Session.GetString(context, IPAddressKey)

	var emptyReservation model.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	templateData := model.TemplateData{
		StringMap:     stringMap,
		FormValidator: form.New(nil),
		Data:          data,
	}

	return &templateData

}

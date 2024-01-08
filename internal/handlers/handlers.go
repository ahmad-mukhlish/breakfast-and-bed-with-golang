package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/helper"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/repository"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/config"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/form"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/model"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/renders"
)

type HandlerRepository struct {
	AppConfig    *config.AppConfig
	DBRepository repository.DatabaseRepository
}

var Repo *HandlerRepository

func CreateRepository(appConfig *config.AppConfig, dbrepository repository.DatabaseRepository) *HandlerRepository {

	return &HandlerRepository{
		AppConfig:    appConfig,
		DBRepository: dbrepository,
	}

}

func CreateHandlers(repository *HandlerRepository) {
	Repo = repository
}

func (m *HandlerRepository) About(w http.ResponseWriter, r *http.Request) {

	var check = m.DBRepository.GetUsers()
	fmt.Print(check)

	initializedTemplate := initiateTemplate()
	err := renders.ServeTemplate(w, r, "about.page.tmpl", initializedTemplate)
	if err != nil {
		helper.CatchServerError(w, err)
		return
	}

}

func (m *HandlerRepository) General(w http.ResponseWriter, r *http.Request) {

	initializedTemplate := initiateTemplate()
	err := renders.ServeTemplate(w, r, "general.page.tmpl", initializedTemplate)
	if err != nil {
		helper.CatchServerError(w, err)
		return
	}
}

func (m *HandlerRepository) Major(w http.ResponseWriter, r *http.Request) {

	initializedTemplate := initiateTemplate()
	err := renders.ServeTemplate(w, r, "major.page.tmpl", initializedTemplate)
	if err != nil {
		helper.CatchServerError(w, err)
		return
	}
}

func (m *HandlerRepository) Contact(w http.ResponseWriter, r *http.Request) {

	initializedTemplate := initiateTemplate()
	err := renders.ServeTemplate(w, r, "contact.page.tmpl", initializedTemplate)
	if err != nil {
		helper.CatchServerError(w, err)
		return
	}
}

func (m *HandlerRepository) Home(w http.ResponseWriter, r *http.Request) {

	initializedTemplate := initiateTemplate()
	err := renders.ServeTemplate(w, r, "home.page.tmpl", initializedTemplate)
	if err != nil {
		helper.CatchServerError(w, err)
		return
	}

}

func (m *HandlerRepository) Reservation(w http.ResponseWriter, r *http.Request) {

	initializedTemplate := initiateTemplate()
	err := renders.ServeTemplate(w, r, "reservation.page.tmpl", initializedTemplate)
	if err != nil {
		helper.CatchServerError(w, err)
		return
	}

}

func (m *HandlerRepository) PostReservation(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		helper.CatchServerError(w, err)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helper.CatchServerError(w, err)
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helper.CatchServerError(w, err)
	}
	roomId, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helper.CatchServerError(w, err)
	}

	reservation := model.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomId:    roomId,
	}

	formValidator := form.NewValidator(r.PostForm)

	formValidator.ValidateLength("first_name", 3)
	formValidator.Required("first_name", "last_name", "phone")
	formValidator.ValidateEmail("email")

	if !formValidator.IsValid() {

		data := make(map[string]interface{})
		data["reservation"] = reservation

		err = renders.ServeTemplate(w, r, "reservation.page.tmpl", &model.TemplateData{
			Data:          data,
			FormValidator: formValidator,
		})
		if err != nil {
			helper.CatchServerError(w, err)
			return
		}

		return

	} else {
		//store the value of reservation in session
		dbErr := m.DBRepository.InsertReservation(reservation)
		if dbErr != nil {
			helper.CatchServerError(w, err)
		}
		m.AppConfig.Session.Put(r.Context(), "reservation", reservation)

		http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
		return
	}
}

func (m *HandlerRepository) ReservationSummary(w http.ResponseWriter, r *http.Request) {

	reservationInterface := m.AppConfig.Session.Pop(r.Context(), "reservation")

	reservation, ok := reservationInterface.(model.Reservation)
	if !ok {
		m.AppConfig.ErrorLog.Println("error parsing data")
		m.AppConfig.Session.Put(r.Context(), "warning", "Please submit your reservation first :)")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	data := make(map[string]interface{})

	data["reservation"] = reservation

	err := renders.ServeTemplate(w, r, "reservation-summary.page.tmpl", &model.TemplateData{
		Data: data,
	})
	if err != nil {
		helper.CatchServerError(w, err)
		return
	}

}

func (m *HandlerRepository) CheckAvailability(w http.ResponseWriter, r *http.Request) {

	initializedTemplate := initiateTemplate()
	err := renders.ServeTemplate(w, r, "check-availability.page.tmpl", initializedTemplate)
	if err != nil {
		helper.CatchServerError(w, err)
		return
	}

}

func (m *HandlerRepository) PostCheckAvailability(w http.ResponseWriter, r *http.Request) {

	arrival := r.Form.Get("start")
	departure := r.Form.Get("end")

	_, err := w.Write([]byte(fmt.Sprintf("Your arrival date is %s, your departure date is %s", arrival, departure)))
	if err != nil {
		helper.CatchServerError(w, err)
		return
	}
}

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func (m *HandlerRepository) CheckAvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	arrival := r.Form.Get("start")
	departure := r.Form.Get("end")
	message := fmt.Sprintf("Your arrival date is %s, your departure date is %s", arrival, departure)

	response := jsonResponse{Ok: true, Message: message}
	output, err := json.MarshalIndent(response, "", "  ")

	if err != nil {
		helper.CatchServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		helper.CatchServerError(w, err)
		return
	}
}

func initiateTemplate() *model.TemplateData {

	var emptyReservation model.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	templateData := model.TemplateData{
		FormValidator: form.NewValidator(nil),
		Data:          data,
	}

	return &templateData

}

func setupLogger() {
	Repo.AppConfig.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	Repo.AppConfig.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	helper.SetHelperAppConfig(Repo.AppConfig)
}

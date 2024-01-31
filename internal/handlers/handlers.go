package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/config"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/form"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/helper"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/model"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/renders"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/repository"
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

	initializedTemplate := initiateTemplate()
	renders.ServeTemplate(w, r, "about.page.tmpl", initializedTemplate)

}

func (m *HandlerRepository) General(w http.ResponseWriter, r *http.Request) {

	initializedTemplate := initiateTemplate()
	renders.ServeTemplate(w, r, "general.page.tmpl", initializedTemplate)

}

func (m *HandlerRepository) Major(w http.ResponseWriter, r *http.Request) {

	initializedTemplate := initiateTemplate()
	renders.ServeTemplate(w, r, "major.page.tmpl", initializedTemplate)

}

func (m *HandlerRepository) Contact(w http.ResponseWriter, r *http.Request) {

	initializedTemplate := initiateTemplate()
	renders.ServeTemplate(w, r, "contact.page.tmpl", initializedTemplate)

}

func (m *HandlerRepository) Home(w http.ResponseWriter, r *http.Request) {

	initializedTemplate := initiateTemplate()
	renders.ServeTemplate(w, r, "home.page.tmpl", initializedTemplate)

}

func (m *HandlerRepository) Reservation(w http.ResponseWriter, r *http.Request) {

	reservation, ok := m.AppConfig.Session.Get(r.Context(), "reservation").(model.Reservation)

	if !ok {
		handleErrorAndRedirect(m, w, r, "Cannot Make A Reservation Right Now")
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = reservation

	room, err := m.DBRepository.GetRoomById(reservation.RoomId)

	if err != nil {
		handleErrorAndRedirect(m, w, r, "Cannot Find The Room")
		return
	}

	data["room_name"] = room.RoomName
	data["arrival"] = reservation.StartDate.Format("Monday, 02 January 2006")
	data["departure"] = reservation.EndDate.Format("Monday, 02 January 2006")
	data["arrivalSQL"] = reservation.StartDate.Format("2006-01-02")
	data["departureSQL"] = reservation.EndDate.Format("2006-01-02")

	templateData := model.TemplateData{
		FormValidator: form.NewValidator(nil),
		Data:          data,
	}

	reservation.Room.RoomName = room.RoomName
	m.AppConfig.Session.Put(r.Context(), "reservation", reservation)

	_ = renders.ServeTemplate(w, r, "reservation.page.tmpl", &templateData)

}

func (m *HandlerRepository) PostReservation(w http.ResponseWriter, r *http.Request) {

	//parse the data
	err := r.ParseForm()
	if err != nil {
		handleErrorAndRedirect(m, w, r, err.Error())
		return
	}

	//get the reservation from session
	reservation, ok := m.AppConfig.Session.Get(r.Context(), "reservation").(model.Reservation)
	if !ok {
		handleErrorAndRedirect(m, w, r, "Cannot get reservation from session")
		return
	}

	//parse the data
	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	//form validator
	formValidator := form.NewValidator(r.PostForm)
	formValidator.ValidateLength("first_name", 3)
	formValidator.Required("first_name", "last_name", "phone")
	formValidator.ValidateEmail("email")

	if !formValidator.IsValid() {

		data := make(map[string]interface{})
		data["reservation"] = reservation

		_ = renders.ServeTemplate(w, r, "reservation.page.tmpl", &model.TemplateData{
			Data:          data,
			FormValidator: formValidator,
		})

		return

	}

	//store the value of reservation in session
	reservationId, dbErr := m.DBRepository.InsertReservation(reservation)
	if dbErr != nil {
		handleErrorAndRedirect(m, w, r, dbErr.Error())
		return
	}

	roomRestriction := model.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomId:        reservation.RoomId,
		ReservationId: reservationId,
		RestrictionId: 1,
	}

	dbErr = m.DBRepository.InsertRoomRestriction(roomRestriction)
	if dbErr != nil {
		handleErrorAndRedirect(m, w, r, dbErr.Error())
		return
	}

	m.AppConfig.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
	return

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
	data["arrival"] = reservation.StartDate.Format("Monday, 02 January 2006")
	data["departure"] = reservation.EndDate.Format("Monday, 02 January 2006")

	renders.ServeTemplate(w, r, "reservation-summary.page.tmpl", &model.TemplateData{
		Data: data,
	})

}

func (m *HandlerRepository) CheckAvailability(w http.ResponseWriter, r *http.Request) {

	initializedTemplate := initiateTemplate()
	renders.ServeTemplate(w, r, "check-availability.page.tmpl", initializedTemplate)

}

func (m *HandlerRepository) PostCheckAvailability(w http.ResponseWriter, r *http.Request) {

	arrival := r.Form.Get("start")
	departure := r.Form.Get("end")
	rooms, err := m.DBRepository.GetAvailableRooms(arrival, departure)

	if err != nil {
		helper.CatchServerError(w, err)
		return
	}

	if len(rooms) > 0 {

		startDate, err := helper.ConvertStringSQLToTime(arrival, "01/02/2006")
		if err != nil {
			helper.CatchServerError(w, err)
			return
		}

		endDate, err := helper.ConvertStringSQLToTime(departure, "01/02/2006")
		if err != nil {
			helper.CatchServerError(w, err)
			return
		}

		reservationWithDates := &model.Reservation{StartDate: startDate, EndDate: endDate}

		m.AppConfig.Session.Put(r.Context(), "reservation",
			reservationWithDates)

		data := make(map[string]interface{})

		data["rooms"] = rooms

		templateWithRoomData := &model.TemplateData{
			Data: data,
		}

		renders.ServeTemplate(w, r, "available-rooms.page.tmpl", templateWithRoomData)

	} else {
		m.AppConfig.Session.Put(r.Context(), "warning", "No Available Rooms :)")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func (m *HandlerRepository) CheckAvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	arrival := r.Form.Get("start")
	departure := r.Form.Get("end")
	roomIdParam := r.Form.Get("room_id")

	roomId, err := strconv.Atoi(roomIdParam)

	if err != nil {
		helper.CatchServerError(w, err)
		return
	}

	var isAvail bool

	isAvail, err = m.DBRepository.CheckAvailabilityForRoomById(arrival, departure, roomId)

	if err != nil {
		helper.CatchServerError(w, err)
		return
	}

	var message string
	if isAvail {
		message = "The room is available"
		startDate, err := helper.ConvertStringSQLToTime(arrival, "01/02/2006")
		if err != nil {
			helper.CatchServerError(w, err)
			return
		}

		endDate, err := helper.ConvertStringSQLToTime(departure, "01/02/2006")
		if err != nil {
			helper.CatchServerError(w, err)
			return
		}

		reservationWithDates := &model.Reservation{StartDate: startDate, EndDate: endDate}
		m.AppConfig.Session.Put(r.Context(), "reservation",
			reservationWithDates)

	} else {
		message = "Sorry, the room is not available"
	}

	var output []byte
	response := jsonResponse{Ok: isAvail, Message: message}
	output, err = json.MarshalIndent(response, "", "  ")

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

func (m *HandlerRepository) CheckRoom(w http.ResponseWriter, r *http.Request) {

	idFromParam := chi.URLParam(r, "id")

	roomId, err := strconv.Atoi(idFromParam)

	if err != nil {

		helper.CatchServerError(w, err)
		return

	}

	reservation, ok := m.AppConfig.Session.Get(r.Context(), "reservation").(model.Reservation)

	if !ok {
		helper.CatchServerError(w, err)
		return
	}

	reservation.RoomId = roomId

	m.AppConfig.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation", http.StatusSeeOther)

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

func handleErrorAndRedirect(m *HandlerRepository, w http.ResponseWriter, r *http.Request, errorMessage string) {
	m.AppConfig.Session.Put(r.Context(), "error", errorMessage)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	helper.CatchServerError(w, errors.New(errorMessage))

}

func setupLogger() {
	Repo.AppConfig.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	Repo.AppConfig.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	helper.SetHelperAppConfig(Repo.AppConfig)
}

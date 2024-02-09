package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
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

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
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
	_ = renders.ServeTemplate(w, r, "about.page.tmpl", initializedTemplate)

}

func (m *HandlerRepository) General(w http.ResponseWriter, r *http.Request) {

	initializedTemplate := initiateTemplate()
	_ = renders.ServeTemplate(w, r, "general.page.tmpl", initializedTemplate)

}

func (m *HandlerRepository) Major(w http.ResponseWriter, r *http.Request) {

	initializedTemplate := initiateTemplate()
	_ = renders.ServeTemplate(w, r, "major.page.tmpl", initializedTemplate)

}

func (m *HandlerRepository) Contact(w http.ResponseWriter, r *http.Request) {

	initializedTemplate := initiateTemplate()
	_ = renders.ServeTemplate(w, r, "contact.page.tmpl", initializedTemplate)

}

func (m *HandlerRepository) Home(w http.ResponseWriter, r *http.Request) {

	initializedTemplate := initiateTemplate()
	_ = renders.ServeTemplate(w, r, "home.page.tmpl", initializedTemplate)

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

	//send email to guest
	sendEmail(m, reservation.Email, "Reservation Notification", buildEmailNotificationGuest(reservation), "basic.html")

	//send email to owner
	sendEmail(m, "owner@breakfast-bed.com", "Owner's Reservation Notification", buildEmailNotificationOwner(reservation), "basic.html")

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

	_ = renders.ServeTemplate(w, r, "reservation-summary.page.tmpl", &model.TemplateData{
		Data: data,
	})

}

func (m *HandlerRepository) CheckAvailability(w http.ResponseWriter, r *http.Request) {

	initializedTemplate := initiateTemplate()
	_ = renders.ServeTemplate(w, r, "check-availability.page.tmpl", initializedTemplate)

}

func (m *HandlerRepository) PostCheckAvailability(w http.ResponseWriter, r *http.Request) {

	errorParse := r.ParseForm()
	if errorParse != nil {
		handleErrorAndRedirect(m, w, r, errorParse.Error())
		return
	}
	arrival := r.Form.Get("start")
	departure := r.Form.Get("end")

	//parse into time and check errors
	startDate, err := helper.ConvertStringSQLToTime(arrival, "02-01-2006")
	if err != nil {
		handleErrorAndRedirect(m, w, r, err.Error())
		return
	}

	endDate, err := helper.ConvertStringSQLToTime(departure, "02-01-2006")
	if err != nil {
		handleErrorAndRedirect(m, w, r, err.Error())
		return
	}

	startDateSQLFormat := startDate.Format("2006-01-02")
	endDateSQLFormat := endDate.Format("2006-01-02")

	rooms, dbErr := m.DBRepository.GetAvailableRooms(
		startDateSQLFormat, endDateSQLFormat)

	if dbErr != nil {
		handleErrorAndRedirect(m, w, r, dbErr.Error())
		return
	}

	//room is not available
	if len(rooms) <= 0 {
		m.AppConfig.Session.Put(r.Context(), "warning", "No Available Rooms :)")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//save to the session
	reservationWithDates := &model.Reservation{StartDate: startDate, EndDate: endDate}
	m.AppConfig.Session.Put(r.Context(), "reservation",
		reservationWithDates)

	//serve to template
	data := make(map[string]interface{})
	data["rooms"] = rooms
	templateWithRoomData := &model.TemplateData{
		Data: data,
	}
	_ = renders.ServeTemplate(w, r, "available-rooms.page.tmpl", templateWithRoomData)

}

func (m *HandlerRepository) PostCheckAvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	errorParse := r.ParseForm()
	if errorParse != nil {
		handleErrorAndRedirect(m, w, r, errorParse.Error())
		return
	}

	roomIdParam := r.Form.Get("room_id")

	roomId, parseErr := strconv.Atoi(roomIdParam)
	if parseErr != nil {
		handleErrorAndRedirect(m, w, r, parseErr.Error())
		return
	}

	arrival := r.Form.Get("start")
	departure := r.Form.Get("end")

	isAvail, dbErr := m.DBRepository.CheckAvailabilityForRoomById(arrival, departure, roomId)
	if dbErr != nil {
		handleErrorAndRedirect(m, w, r, dbErr.Error())
		return
	}

	var message string
	if isAvail {
		message = "The room is available"
		startDate, err := helper.ConvertStringSQLToTime(arrival, "02-01-2006")
		if err != nil {
			handleErrorAndRedirect(m, w, r, err.Error())
			return
		}

		endDate, err := helper.ConvertStringSQLToTime(departure, "02-01-2006")
		if err != nil {
			handleErrorAndRedirect(m, w, r, err.Error())
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
	output, _ = json.MarshalIndent(response, "", "  ")

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(output)

}

func (m *HandlerRepository) CheckRoom(w http.ResponseWriter, r *http.Request) {

	idFromParam := chi.URLParam(r, "id")

	roomId, err := strconv.Atoi(idFromParam)
	if err != nil {
		handleErrorAndRedirect(m, w, r, err.Error())
		return
	}

	reservation, ok := m.AppConfig.Session.Get(r.Context(), "reservation").(model.Reservation)
	if !ok {
		handleErrorAndRedirect(m, w, r, "Cannot get reservation from session")
		return
	}

	reservation.RoomId = roomId
	m.AppConfig.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation", http.StatusSeeOther)

}

func (m *HandlerRepository) Login(w http.ResponseWriter, r *http.Request) {

	initializedTemplate := initiateTemplate()
	_ = renders.ServeTemplate(w, r, "login.page.tmpl", initializedTemplate)

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

func sendEmail(m *HandlerRepository, emailTo, subject, emailContent, templateName string) {

	mail := model.MailData{
		To:           emailTo,
		From:         "admin@breakfast-and-bed.com",
		Content:      emailContent,
		Subject:      subject,
		TemplateName: templateName,
	}

	m.AppConfig.MailChan <- mail
}

func buildEmailNotificationGuest(reservation model.Reservation) string {

	emailTemplate := `<strong> Reservation Confirmation </strong> <br>
                      <br> Dear %s %s, <br> <br>
                      This is the confirmation of your reservation<br> <br> The date will be from %s to %s <br> in our breakfast and bed hotel <br><br>
                      Your room name will be %s <br><br>

                      Best regards, <br>
                      Breakfast and Bed.
                     
`

	return fmt.Sprintf(emailTemplate, reservation.FirstName, reservation.LastName,
		reservation.StartDate.Format("02-01-2006"), reservation.EndDate.Format("02-01-2006"),
		reservation.Room.RoomName)
}

func buildEmailNotificationOwner(reservation model.Reservation) string {

	emailTemplate := `<strong> Owner's Reservation Confirmation </strong> <br>
                      <br> Dear Owner <br> <br>
                      This is the confirmation of the reservation<br> <br> The date will be from %s to %s <br> in our breakfast and bed hotel <br><br>
                      The room name will be %s <br><br>

                      The name of our guest is %s %s <br>
                      (Phone Number : %s) <br> <br>

                      Best regards, <br>
                      Breakfast and Bed Admin.
                     
`

	return fmt.Sprintf(emailTemplate,
		reservation.StartDate.Format("02-01-2006"), reservation.EndDate.Format("02-01-2006"),
		reservation.Room.RoomName, reservation.FirstName, reservation.LastName, reservation.Phone)
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

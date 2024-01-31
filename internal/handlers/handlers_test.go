package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"testing"

	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/model"
)

type postData struct {
	key   string
	value string
}

type test struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}

var tests = []test{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"major", "/major", "GET", []postData{}, http.StatusOK},
	{"general", "/general", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"reservation-summary", "/reservation-summary",
		"GET", []postData{}, http.StatusOK},
	{"check-availability", "/check-availability",
		"GET", []postData{}, http.StatusOK},
	{"no route", "/none", "GET", []postData{}, http.StatusNotFound},
	//POST

}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, currentTest := range tests {
		if currentTest.method == "GET" {

			mockedClient := testServer.Client()
			mockedBasePath := testServer.URL
			mockedFullURl := mockedBasePath + currentTest.url

			response, err := mockedClient.Get(mockedFullURl)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if response.StatusCode != currentTest.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d",
					currentTest.name, currentTest.expectedStatusCode, response.StatusCode)
			}

		} else if currentTest.method == "POST" {

			mockedExpectedParam := url.Values{}
			for _, param := range currentTest.params {
				mockedExpectedParam.Add(param.key, param.value)
			}

			mockedClient := testServer.Client()
			mockedBasePath := testServer.URL
			mockedFullURl := mockedBasePath + currentTest.url

			response, err := mockedClient.PostForm(mockedFullURl, mockedExpectedParam)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if response.StatusCode != currentTest.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d",
					currentTest.name, currentTest.expectedStatusCode, response.StatusCode)
			}

		}
	}
}

func Test_Reservation_Success(t *testing.T) {

	// create a request instance
	request, _ := http.NewRequest("GET", "/reservation", nil)

	// get the context of the request
	currentContext := request.Context()

	//make the context session-able
	sessionedContext, err := AppConfig.Session.Load(currentContext, request.Header.Get("X-Session"))
	if err != nil {
		log.Fatal(err)
	}

	//make the request session-able
	request = request.WithContext(sessionedContext)

	// put the reservation into the session via sessioned context
	reservation := model.Reservation{
		RoomId: 1,
		Room: model.Room{
			Id:       1,
			RoomName: "Generals",
		},
	}

	AppConfig.Session.Put(sessionedContext, "reservation", reservation)

	rr := httptest.NewRecorder()
	mockedHandler := http.HandlerFunc(Repo.Reservation)
	mockedHandler.ServeHTTP(rr, request)

	if rr.Code != http.StatusOK {
		t.Error("Status Code Not OK")
	}

}

func Test_Reservation_NoSession(t *testing.T) {

	// create a request instance
	request, _ := http.NewRequest("GET", "/reservation", nil)

	// get the context of the request
	currentContext := request.Context()

	//make the context session-able
	sessionedContext, err := AppConfig.Session.Load(currentContext, request.Header.Get("X-Session"))
	if err != nil {
		log.Fatal(err)
	}

	//make the request session-able
	request = request.WithContext(sessionedContext)

	rr := httptest.NewRecorder()
	mockedHandler := http.HandlerFunc(Repo.Reservation)
	mockedHandler.ServeHTTP(rr, request)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Error("Status Code Must Be Temporary Redirect")
	}

}

func Test_Reservation_NoRoomInDB(t *testing.T) {

	// create a request instance
	request, _ := http.NewRequest("GET", "/reservation", nil)

	// get the context of the request
	currentContext := request.Context()

	//make the context session-able
	sessionedContext, err := AppConfig.Session.Load(currentContext, request.Header.Get("X-Session"))
	if err != nil {
		log.Fatal(err)
	}

	//make the request session-able
	request = request.WithContext(sessionedContext)

	// put the reservation into the session via sessioned context
	reservation := model.Reservation{
		RoomId: 3,
		Room: model.Room{
			Id:       3,
			RoomName: "Nil",
		},
	}

	AppConfig.Session.Put(sessionedContext, "reservation", reservation)

	rr := httptest.NewRecorder()
	mockedHandler := http.HandlerFunc(Repo.Reservation)
	mockedHandler.ServeHTTP(rr, request)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Status Code Must Be Temporary Redirect but instead get %d", rr.Code)
	}
}

func Test_PostReservation_Success(t *testing.T) {

	reqBody := url.Values{}

	reqBody.Add("first_name", "Ahmad")
	reqBody.Add("last_name", "Mukhlis")
	reqBody.Add("email", "ahmad@mukhlis.com")
	reqBody.Add("phone", "11111111111")

	// create a request instance
	request, _ := http.NewRequest("POST", "/reservation", strings.NewReader(reqBody.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// get the context of the request
	currentContext := request.Context()

	//make the context session-able
	sessionedContext, err := AppConfig.Session.Load(currentContext, request.Header.Get("X-Session"))
	if err != nil {
		log.Fatal(err)
	}

	//make the request session-able
	request = request.WithContext(sessionedContext)

	// put the reservation into the session via sessioned context
	reservation := model.Reservation{
		RoomId: 1,
		Room: model.Room{
			Id:       1,
			RoomName: "Generals",
		},
	}

	AppConfig.Session.Put(sessionedContext, "reservation", reservation)

	rr := httptest.NewRecorder()
	mockedHandler := http.HandlerFunc(Repo.PostReservation)
	mockedHandler.ServeHTTP(rr, request)

	if rr.Code != http.StatusSeeOther {
		t.Error("Status Code Not See Other")
	}

}

func Test_PostReservation_NoSession(t *testing.T) {

	reqBody := url.Values{}

	reqBody.Add("first_name", "Ahmad")
	reqBody.Add("last_name", "Mukhlis")
	reqBody.Add("email", "ahmad@mukhlis.com")
	reqBody.Add("phone", "11111111111")

	// create a request instance
	request, _ := http.NewRequest("POST", "/reservation", strings.NewReader(reqBody.Encode()))

	// get the context of the request
	currentContext := request.Context()

	//make the context session-able
	sessionedContext, err := AppConfig.Session.Load(currentContext, request.Header.Get("X-Session"))
	if err != nil {
		log.Fatal(err)
	}

	//make the request session-able
	request = request.WithContext(sessionedContext)

	rr := httptest.NewRecorder()
	mockedHandler := http.HandlerFunc(Repo.PostReservation)
	mockedHandler.ServeHTTP(rr, request)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Error("Status Code Must Be Temporary Redirect")
	}

}

func Test_PostReservation_NoBodyCannotParse(t *testing.T) {

	// create a request instance
	request, _ := http.NewRequest("POST", "/reservation", nil)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// get the context of the request
	currentContext := request.Context()

	//make the context session-able
	sessionedContext, err := AppConfig.Session.Load(currentContext, request.Header.Get("X-Session"))
	if err != nil {
		log.Fatal(err)
	}

	//make the request session-able
	request = request.WithContext(sessionedContext)

	rr := httptest.NewRecorder()
	mockedHandler := http.HandlerFunc(Repo.PostReservation)
	mockedHandler.ServeHTTP(rr, request)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Error("Status Code Must Be Temporary Redirect")
	}

}

func Test_PostReservation_InvalidData(t *testing.T) {

	reqBody := url.Values{}

	reqBody.Add("last_name", "Mukhlis")
	reqBody.Add("email", "ahmad@mukhlis.com")
	reqBody.Add("phone", "11111111111")

	// create a request instance
	request, _ := http.NewRequest("POST", "/reservation", strings.NewReader(reqBody.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// get the context of the request
	currentContext := request.Context()

	//make the context session-able
	sessionedContext, err := AppConfig.Session.Load(currentContext, request.Header.Get("X-Session"))
	if err != nil {
		log.Fatal(err)
	}

	//make the request session-able
	request = request.WithContext(sessionedContext)

	// put the reservation into the session via sessioned context
	reservation := model.Reservation{
		RoomId: 1,
		Room: model.Room{
			Id:       1,
			RoomName: "Generals",
		},
	}

	AppConfig.Session.Put(sessionedContext, "reservation", reservation)

	rr := httptest.NewRecorder()
	mockedHandler := http.HandlerFunc(Repo.PostReservation)
	mockedHandler.ServeHTTP(rr, request)

	if rr.Code != http.StatusOK {
		t.Error("Status Code Not OK")
	}

}

func Test_PostReservation_CannotInsertReservation(t *testing.T) {

	reqBody := url.Values{}

	reqBody.Add("first_name", "Ahmad")
	reqBody.Add("last_name", "Mukhlis")
	reqBody.Add("email", "ahmad@mukhlis.com")
	reqBody.Add("phone", "11111111111")

	// create a request instance
	request, _ := http.NewRequest("POST", "/reservation", strings.NewReader(reqBody.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// get the context of the request
	currentContext := request.Context()

	//make the context session-able
	sessionedContext, err := AppConfig.Session.Load(currentContext, request.Header.Get("X-Session"))
	if err != nil {
		log.Fatal(err)
	}

	//make the request session-able
	request = request.WithContext(sessionedContext)

	// put the reservation into the session via sessioned context
	reservation := model.Reservation{
		RoomId: 1000,
		Room: model.Room{
			Id:       1000,
			RoomName: "Generals",
		},
	}

	AppConfig.Session.Put(sessionedContext, "reservation", reservation)

	rr := httptest.NewRecorder()
	mockedHandler := http.HandlerFunc(Repo.PostReservation)
	mockedHandler.ServeHTTP(rr, request)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Status Code Not Temporary Redirect and Get %d", rr.Code)
	}

}

func Test_PostReservation_CannotInsertRoomRestriction(t *testing.T) {

	reqBody := url.Values{}

	reqBody.Add("first_name", "Ahmad")
	reqBody.Add("last_name", "Mukhlis")
	reqBody.Add("email", "ahmad@mukhlis.com")
	reqBody.Add("phone", "11111111111")

	// create a request instance
	request, _ := http.NewRequest("POST", "/reservation", strings.NewReader(reqBody.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// get the context of the request
	currentContext := request.Context()

	//make the context session-able
	sessionedContext, err := AppConfig.Session.Load(currentContext, request.Header.Get("X-Session"))
	if err != nil {
		log.Fatal(err)
	}

	//make the request session-able
	request = request.WithContext(sessionedContext)

	// put the reservation into the session via sessioned context
	reservation := model.Reservation{
		RoomId: 10000,
		Room: model.Room{
			Id:       10000,
			RoomName: "Generals",
		},
	}

	AppConfig.Session.Put(sessionedContext, "reservation", reservation)

	rr := httptest.NewRecorder()
	mockedHandler := http.HandlerFunc(Repo.PostReservation)
	mockedHandler.ServeHTTP(rr, request)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Status Code Not Temporary Redirect and Get %d", rr.Code)
	}

}

func Test_PostCheckAvailability_Success(t *testing.T) {

	reqBody := url.Values{}

	reqBody.Add("start", "01-02-2024")
	reqBody.Add("end", "10-02-2024")
	reqBody.Add("first_name", "Ahmad")
	reqBody.Add("last_name", "Mukhlis")
	reqBody.Add("email", "ahmad@mukhlis.com")
	reqBody.Add("phone", "11111111111")

	// create a request instance
	request, _ := http.NewRequest("POST", "/check-availability", strings.NewReader(reqBody.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// get the context of the request
	currentContext := request.Context()

	//make the context session-able
	sessionedContext, err := AppConfig.Session.Load(currentContext, request.Header.Get("X-Session"))
	if err != nil {
		log.Fatal(err)
	}

	//make the request session-able
	request = request.WithContext(sessionedContext)

	// put the reservation into the session via sessioned context
	reservation := model.Reservation{
		RoomId: 1,
		Room: model.Room{
			Id:       1,
			RoomName: "Generals",
		},
	}

	AppConfig.Session.Put(sessionedContext, "reservation", reservation)

	rr := httptest.NewRecorder()
	mockedHandler := http.HandlerFunc(Repo.PostCheckAvailability)
	mockedHandler.ServeHTTP(rr, request)

	if rr.Code != http.StatusOK {
		t.Error("Status Code Not OK")
	}

}

func Test_PostCheckAvailability_NoBodyCannotParse(t *testing.T) {

	// create a request instance
	request, _ := http.NewRequest("POST", "/check-availability", nil)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// get the context of the request
	currentContext := request.Context()

	//make the context session-able
	sessionedContext, err := AppConfig.Session.Load(currentContext, request.Header.Get("X-Session"))
	if err != nil {
		log.Fatal(err)
	}

	//make the request session-able
	request = request.WithContext(sessionedContext)

	rr := httptest.NewRecorder()
	mockedHandler := http.HandlerFunc(Repo.PostCheckAvailability)
	mockedHandler.ServeHTTP(rr, request)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Error("Status Code Must Be Temporary Redirect")
	}

}

func Test_PostCheckAvailability_AvailableRoomDBError(t *testing.T) {

	reqBody := url.Values{}

	reqBody.Add("start", "02-01-2006")
	reqBody.Add("end", "04-01-2006")
	reqBody.Add("first_name", "Ahmad")
	reqBody.Add("last_name", "Mukhlis")
	reqBody.Add("email", "ahmad@mukhlis.com")
	reqBody.Add("phone", "11111111111")

	// create a request instance
	request, _ := http.NewRequest("POST", "/check-availability", strings.NewReader(reqBody.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// get the context of the request
	currentContext := request.Context()

	//make the context session-able
	sessionedContext, err := AppConfig.Session.Load(currentContext, request.Header.Get("X-Session"))
	if err != nil {
		log.Fatal(err)
	}

	//make the request session-able
	request = request.WithContext(sessionedContext)

	rr := httptest.NewRecorder()
	mockedHandler := http.HandlerFunc(Repo.PostCheckAvailability)
	mockedHandler.ServeHTTP(rr, request)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Error("Status Code Must Be Temporary Redirect")
	}

}

func Test_PostCheckAvailability_NoAvailableRoom(t *testing.T) {

	reqBody := url.Values{}

	reqBody.Add("start", "03-01-2006")
	reqBody.Add("end", "04-01-2006")
	reqBody.Add("first_name", "Ahmad")
	reqBody.Add("last_name", "Mukhlis")
	reqBody.Add("email", "ahmad@mukhlis.com")
	reqBody.Add("phone", "11111111111")

	// create a request instance
	request, _ := http.NewRequest("POST", "/check-availability", strings.NewReader(reqBody.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// get the context of the request
	currentContext := request.Context()

	//make the context session-able
	sessionedContext, err := AppConfig.Session.Load(currentContext, request.Header.Get("X-Session"))
	if err != nil {
		log.Fatal(err)
	}

	//make the request session-able
	request = request.WithContext(sessionedContext)

	rr := httptest.NewRecorder()
	mockedHandler := http.HandlerFunc(Repo.PostCheckAvailability)
	mockedHandler.ServeHTTP(rr, request)

	if rr.Code != http.StatusSeeOther {
		t.Error("Status Code Must Be See Other")
	}

}

func Test_PostCheckAvailability_NoStartDate(t *testing.T) {

	reqBody := url.Values{}

	reqBody.Add("end", "04-01-2006")
	reqBody.Add("first_name", "Ahmad")
	reqBody.Add("last_name", "Mukhlis")
	reqBody.Add("email", "ahmad@mukhlis.com")
	reqBody.Add("phone", "11111111111")

	// create a request instance
	request, _ := http.NewRequest("POST", "/check-availability", strings.NewReader(reqBody.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// get the context of the request
	currentContext := request.Context()

	//make the context session-able
	sessionedContext, err := AppConfig.Session.Load(currentContext, request.Header.Get("X-Session"))
	if err != nil {
		log.Fatal(err)
	}

	//make the request session-able
	request = request.WithContext(sessionedContext)

	rr := httptest.NewRecorder()
	mockedHandler := http.HandlerFunc(Repo.PostCheckAvailability)
	mockedHandler.ServeHTTP(rr, request)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Error("Status Code Must Be Temporary Redirect")
	}

}

func Test_PostCheckAvailability_NoEndDate(t *testing.T) {

	reqBody := url.Values{}

	reqBody.Add("start", "01-04-2050")
	reqBody.Add("first_name", "Ahmad")
	reqBody.Add("last_name", "Mukhlis")
	reqBody.Add("email", "ahmad@mukhlis.com")
	reqBody.Add("phone", "11111111111")

	// create a request instance
	request, _ := http.NewRequest("POST", "/check-availability", strings.NewReader(reqBody.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// get the context of the request
	currentContext := request.Context()

	//make the context session-able
	sessionedContext, err := AppConfig.Session.Load(currentContext, request.Header.Get("X-Session"))
	if err != nil {
		log.Fatal(err)
	}

	//make the request session-able
	request = request.WithContext(sessionedContext)

	rr := httptest.NewRecorder()
	mockedHandler := http.HandlerFunc(Repo.PostCheckAvailability)
	mockedHandler.ServeHTTP(rr, request)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Status Code Must Be Temporary Redirect Got %d", rr.Code)
	}

}

func Test_CheckRoom_Success(t *testing.T) {

	// create a request instance
	request, _ := http.NewRequest("GET", "/check/room/{id}", nil)

	// get the context of the request
	currentContext := request.Context()

	//make the context session-able
	sessionedContext, err := AppConfig.Session.Load(currentContext, request.Header.Get("X-Session"))
	if err != nil {
		log.Fatal(err)
	}

	//make the request session-able
	request = request.WithContext(sessionedContext)

	// set the RequestURI on the request so that we can grab the ID
	// from the URL

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")

	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

	// put the reservation into the session via sessioned context
	reservation := model.Reservation{
		RoomId: 1,
		Room: model.Room{
			Id:       1,
			RoomName: "Generals",
		},
	}

	AppConfig.Session.Put(sessionedContext, "reservation", reservation)

	rr := httptest.NewRecorder()
	mockedHandler := http.HandlerFunc(Repo.CheckRoom)
	mockedHandler.ServeHTTP(rr, request)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Status Code Not See Other")
	}

}

func Test_CheckRoom_NoID(t *testing.T) {

	// create a request instance
	request, _ := http.NewRequest("GET", "/check/room/{id}", nil)

	// get the context of the request
	currentContext := request.Context()

	//make the context session-able
	sessionedContext, err := AppConfig.Session.Load(currentContext, request.Header.Get("X-Session"))
	if err != nil {
		log.Fatal(err)
	}

	//make the request session-able
	request = request.WithContext(sessionedContext)

	// set the RequestURI on the request so that we can grab the ID
	// from the URL

	// put the reservation into the session via sessioned context
	reservation := model.Reservation{
		RoomId: 1,
		Room: model.Room{
			Id:       1,
			RoomName: "Generals",
		},
	}

	AppConfig.Session.Put(sessionedContext, "reservation", reservation)

	rr := httptest.NewRecorder()
	mockedHandler := http.HandlerFunc(Repo.CheckRoom)
	mockedHandler.ServeHTTP(rr, request)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Status Code Not Temporary Redirect")
	}

}

func Test_CheckRoom_NoSession(t *testing.T) {

	// create a request instance
	request, _ := http.NewRequest("GET", "/check/room/{id}", nil)

	// get the context of the request
	currentContext := request.Context()

	//make the context session-able
	sessionedContext, err := AppConfig.Session.Load(currentContext, request.Header.Get("X-Session"))
	if err != nil {
		log.Fatal(err)
	}

	//make the request session-able
	request = request.WithContext(sessionedContext)

	// set the RequestURI on the request so that we can grab the ID
	// from the URL

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")

	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

	AppConfig.Session.Put(sessionedContext, "reservation", nil)

	rr := httptest.NewRecorder()
	mockedHandler := http.HandlerFunc(Repo.CheckRoom)
	mockedHandler.ServeHTTP(rr, request)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Status Code Not Temporary Redirect")
	}

}

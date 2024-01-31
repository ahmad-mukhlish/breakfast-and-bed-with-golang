package handlers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"

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

func createDummyPostAvailabilityParams() []postData {

	return []postData{
		{"start", "2022-12-11"},
		{"end", "2022-12-12"},
	}

}

func createDummyPostReservationParams() []postData {

	return []postData{
		{"first_name", "John"},
		{"last_name", "Doe"},
		{"email", "john@doe.com"},
		{"phone", "555-555"},
	}

}

func createDummyPostReservationInvalidParams() []postData {

	return []postData{
		{"first_name", ""},
		{"last_name", "Doe"},
		{"email", "john@doe.com"},
		{"phone", "555-555"},
	}

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

func Test_Reservation(t *testing.T) {

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

func Test_ReservationNoSession(t *testing.T) {

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

func Test_ReservationNoRoomInDB(t *testing.T) {

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

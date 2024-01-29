package handlers

import (
	"log"
	"net/http"
	"net/http/httptest"

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

func Test_Reservation(t *testing.T) {

	// create a request instance
	request, _ := http.NewRequest("GET", "/make-reservation", nil)

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
	request, _ := http.NewRequest("GET", "/make-reservation", nil)

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
	request, _ := http.NewRequest("GET", "/make-reservation", nil)

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
		t.Error("Status Code Must Be Temporary Redirect")
	}
}

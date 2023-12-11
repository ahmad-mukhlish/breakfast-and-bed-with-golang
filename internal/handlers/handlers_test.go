package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
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
	{"reservation", "/reservation", "GET", []postData{}, http.StatusOK},
	{"reservation-summary", "/reservation-summary",
		"GET", []postData{}, http.StatusOK},
	{"check-availability", "/check-availability",
		"GET", []postData{}, http.StatusOK},
	{"no route", "/none", "GET", []postData{}, http.StatusNotFound},
	//POST
	{"post check availability", "/check-availability",
		"POST", createDummyPostAvailabilityParams(), http.StatusOK},
	{"post check availability json", "/check-availability/json",
		"POST", createDummyPostAvailabilityParams(), http.StatusOK},
	{"post reservation", "/reservation",
		"POST", createDummyPostReservationParams(), http.StatusOK},
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

package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type postData map[string]string

type test struct {
	name               string
	url                string
	method             string
	params             postData
	expectedStatusCode int
}

var tests = []test{
	{"home", "/", "GET", make(map[string]string), http.StatusOK},
	{"about", "/about", "GET", make(map[string]string), http.StatusOK},
	{"major", "/major", "GET", make(map[string]string), http.StatusOK},
	{"general", "/general", "GET", make(map[string]string), http.StatusOK},
	{"contact", "/contact", "GET", make(map[string]string), http.StatusOK},
	{"reservation", "/reservation", "GET", make(map[string]string), http.StatusOK},
	{"reservation-summary", "/reservation-summary", "GET", make(map[string]string), http.StatusOK},
	{"check-availability", "/check-availability", "GET", make(map[string]string), http.StatusOK},
	{"no route", "/none", "GET", make(map[string]string), http.StatusNotFound},
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

		}
	}
}

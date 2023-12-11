package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {

	mockedHandler := mockHandler{}
	handler := NoSurf(mockedHandler)

	_, ok := handler.(http.Handler)

	if !ok {
		t.Error("type is not http.Handler")
	}
}

func TestSessionLoad(t *testing.T) {

	mockedHandler := mockHandler{}
	handler := SessionLoad(mockedHandler)

	_, ok := handler.(http.Handler)

	if !ok {
		t.Error("type is not http.Handler")
	}
}

func TestCreateCookie(t *testing.T) {

	cookie := CreateCookie()

	ok := isCookieType(cookie)

	if !ok {
		t.Error("type is not http Cookie")
	}
}

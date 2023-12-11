package main

import (
	"log"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

type mockHandler struct {
}

func (mockHandler mockHandler) ServeHTTP(http.ResponseWriter, *http.Request) {

}

func isCookieType(testedValue interface{}) bool {
	switch v := testedValue.(type) {
	case http.Cookie:
		return true
	default:
		log.Fatal(v)
		return false

	}
}

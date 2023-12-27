package helper

import (
	"net/http/httptest"
	"testing"
)

func TestSetHelperAppConfig(t *testing.T) {

	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()
}

func TestCatchServerError(t *testing.T) {

	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()
}

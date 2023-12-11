package main

import (
	"github.com/go-chi/chi/v5"
	"testing"
)

func TestHandleRoute(t *testing.T) {

	AppConfig, _ = setupConfig()
	handledRoute := HandleRoute()

	_, ok := handledRoute.(*chi.Mux)

	if !ok {
		t.Error("type is not *chi.Mux")
	}
}

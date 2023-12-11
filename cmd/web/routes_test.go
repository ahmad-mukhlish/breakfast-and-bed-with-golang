package main

import (
	"net/http"
	"testing"
)

func TestHandleRoute(t *testing.T) {
	handledRoute := HandleRoute()

	_, ok := handledRoute.(http.Handler)

	if !ok {
		t.Error("type is not http.Handler")
	}
}

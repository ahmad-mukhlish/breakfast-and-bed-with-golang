package main

import (
	"testing"
)

func TestSetupServer(t *testing.T) {
	err := setupServer()
	if err != nil {
		t.Error("unexpected error:", err)
	}

}

func TestStartServer(t *testing.T) {

	//TODO @ahmad-mukhlis fix this
	//err := setupServer()
	//srv, err := startServer(":8080")
	//err = srv.Shutdown(context.Background())
	//if err != nil {
	//	return
	//}
	//
	//if err != nil {
	//	t.Error("unexpected error:", err)
	//}

}

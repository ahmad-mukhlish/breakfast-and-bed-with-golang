package main

import (
	"context"
	"testing"
	"time"
)

func TestSetupServer(t *testing.T) {
	err := setupServer()
	if err != nil {
		t.Error("unexpected error:", err)
	}

}

func TestStartServer(t *testing.T) {

	err := setupServer()
	if err != nil {
		t.Error("unexpected error:", err)
	}

	srv, err := startServer(":45566")
	go func() {
		time.Sleep(1 * time.Second)
		err = srv.Shutdown(context.Background())
		if err != nil {
			return
		}
	}()

	if err != nil {
		t.Error("unexpected error:", err)
	}

}

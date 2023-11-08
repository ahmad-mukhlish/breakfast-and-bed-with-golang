package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

func printToConsole() {
	fmt.Println("This is middleware in action")
}

func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			printToConsole()
			next.ServeHTTP(w, r)
		})
}

func createCookie() http.Cookie {
	return http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
}

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	cookie := createCookie()

	csrfHandler.SetBaseCookie(cookie)
	return csrfHandler

}

func SessionLoad(next http.Handler) http.Handler {
	return AppConfig.Session.LoadAndSave(next)
}

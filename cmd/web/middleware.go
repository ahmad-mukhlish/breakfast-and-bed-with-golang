package main

import (
	"github.com/justinas/nosurf"
	"net/http"
)

func CreateCookie() http.Cookie {
	return http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
}

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	cookie := CreateCookie()

	csrfHandler.SetBaseCookie(cookie)
	return csrfHandler

}

func SessionLoad(next http.Handler) http.Handler {
	return AppConfig.Session.LoadAndSave(next)
}

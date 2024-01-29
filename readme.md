# Breakdast and Bed with GoLang

This is the repository for Breakfast and Bed With GoLang

How To Test With Coverage

go test -coverprofile=coverage.out
go tool cover -html=coverage.out


- Built in Go version 1.21
- Uses the [chi router](github.com/go-chi/chi/v5)
- Uses the [alex edwards SCS](github.com/alexedwards/scs/v2)
- Uses the [nosurf](github.com/justinas/nosurf)
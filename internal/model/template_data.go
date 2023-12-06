package model

import "github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/form"

type TemplateData struct {
	StringMap     map[string]string
	Data          map[string]interface{}
	FormValidator *form.FormValidator
	Error         string
	Flash         string
	Warning       string
}

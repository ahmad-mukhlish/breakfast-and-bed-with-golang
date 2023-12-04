package model

import "github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/form"

type TemplateData struct {
	StringMap     map[string]string
	FormValidator *form.FormValidator
}

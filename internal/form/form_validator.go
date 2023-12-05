package form

import (
	"net/http"
	"net/url"
)

type FormValidator struct {
	Data      url.Values
	FormError FormError
}

func New(data url.Values) *FormValidator {
	return &FormValidator{
		data,
		FormError(map[string][]string{}),
	}
}

func (f *FormValidator) Has(requiredField string, r *http.Request) bool {

	result := r.Form.Get(requiredField)
	isRequiredResultFilled := result != ""

	if !isRequiredResultFilled {
		f.FormError.Add(requiredField, "This field is mandatory")
	}

	return isRequiredResultFilled

}

func (f *FormValidator) IsValid() bool {
	return len(f.FormError) == 0
}

func (f *FormValidator) Required(fields ...string) {

	for _, field := range fields {

		if f.Data.Get(field) == "" {
			f.FormError.Add(field, "This field is Required")
		}

	}

}

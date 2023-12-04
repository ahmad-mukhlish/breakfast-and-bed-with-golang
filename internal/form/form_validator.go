package form

import (
	"net/http"
	"net/url"
)

type FormValidator struct {
	data      url.Values
	formError FormError
}

func New(data url.Values) *FormValidator {
	return &FormValidator{
		data,
		FormError(map[string][]string{}),
	}
}

func HasRequiredField(requiredField string, r http.Request) bool {

	result := r.Form.Get(requiredField)
	isRequiredResultFilled := result != ""

	return isRequiredResultFilled

}

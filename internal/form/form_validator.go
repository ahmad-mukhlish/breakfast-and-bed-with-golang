package form

import (
	"fmt"
	"net/url"

	"github.com/asaskevich/govalidator"
)

type FormValidator struct {
	Data      url.Values
	FormError FormError
}

func NewValidator(data url.Values) *FormValidator {
	return &FormValidator{
		data,
		map[string][]string{},
	}
}

func (f *FormValidator) Has(requiredField string) bool {

	result := f.Data.Get(requiredField)
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

func (f *FormValidator) ValidateLength(field string, minimumLength int) bool {

	if len(f.Data.Get(field)) < minimumLength {
		f.FormError.Add(field, fmt.Sprintf("The minimum length is %d", minimumLength))
		return false
	}
	return true

}

func (f *FormValidator) ValidateEmail(field string) bool {

	validEmail := govalidator.IsEmail(f.Data.Get(field))

	if !validEmail {
		f.FormError.Add(field, "This is Not A Valid Email")
	}

	return validEmail

}

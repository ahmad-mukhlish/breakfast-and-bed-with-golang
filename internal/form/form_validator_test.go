package form

import (
	"net/url"
	"testing"
)

func TestFormValidator_IsValid(t *testing.T) {
	postForm := url.Values{}

	validator := NewValidator(postForm)

	isValid := validator.IsValid()

	if !isValid {
		t.Errorf("Should not return error")
	}

}

func TestFormValidator_Required(t *testing.T) {

	postForm := url.Values{}

	postForm.Add("a", "a")
	postForm.Add("b", "a")
	postForm.Add("c", "a")

	validatorErrorCase := NewValidator(postForm)

	validatorErrorCase.Required("a", "zzz")

	errorCase := validatorErrorCase.IsValid()

	if errorCase {
		t.Errorf("Should not passed, error expected")
	}

	errorMessage := validatorErrorCase.FormError.GetFirstErrorMessage("zzz")

	if errorMessage == "" {
		t.Errorf("Expected error message is filled, because zzz is required but not filled")
	}

	validatorPassedCase := NewValidator(postForm)

	validatorPassedCase.Required("a", "b", "c")

	passedCase := validatorPassedCase.IsValid()

	if !passedCase {
		t.Errorf("Should be passed")
	}

	errorMessageEmpty := validatorPassedCase.FormError.GetFirstErrorMessage("a")

	if errorMessageEmpty != "" {
		t.Errorf("Expected error message is none, because the form should be valid")
	}

}

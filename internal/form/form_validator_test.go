package form

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestFormValidator_IsValid(t *testing.T) {
	mockedReq := httptest.NewRequest("POST", "/testo", nil)
	form := mockedReq.PostForm

	validator := NewValidator(form)

	isValid := validator.IsValid()

	if !isValid {
		t.Errorf("Should not return error")
	}

}

func TestFormValidator_Required(t *testing.T) {

	mockedReq := httptest.NewRequest("POST", "/testo", nil)

	postForm := url.Values{}

	postForm.Add("a", "a")
	postForm.Add("b", "a")
	postForm.Add("c", "a")

	mockedReq.PostForm = postForm

	validatorErrorCase := NewValidator(mockedReq.PostForm)

	validatorErrorCase.Required("a", "zzz")

	errorCase := validatorErrorCase.IsValid()

	if errorCase {
		t.Errorf("Should not passed, error expected")
	}

	validatorPassedCase := NewValidator(mockedReq.PostForm)

	validatorPassedCase.Required("a", "b", "c")

	passedCase := validatorPassedCase.IsValid()

	if !passedCase {
		t.Errorf("Should be passed")
	}

}

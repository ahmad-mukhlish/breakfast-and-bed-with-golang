package form

import (
	"net/http/httptest"
	"testing"
)

func TestFormValidator_IsValid(t *testing.T) {
	mockedReq := httptest.NewRequest("POST", "/testo", nil)
	form := mockedReq.Form

	validator := NewValidator(form)

	isValid := validator.IsValid()

	if !isValid {
		t.Errorf("Should not return error")
	}

}

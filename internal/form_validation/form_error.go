package form_validation

type FormError map[string][]string

func (formError FormError) Add(formField, message string) {
	formError[formField] = append(formError[formField], message)
}

func (formError FormError) GetFirstErrorMessage(formField string) string {
	errorMessages := formError[formField]
	if len(errorMessages) == 0 {
		return ""
	}
	return errorMessages[0]
}

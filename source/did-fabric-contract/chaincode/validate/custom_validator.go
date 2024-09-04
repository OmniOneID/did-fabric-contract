// Copyright 2024 Raonsecure

package validate

import (
	"gopkg.in/go-playground/validator.v9"
)

func isParentExist(fl validator.FieldLevel) bool {
	return !fl.Parent().IsZero()
}

func optionalRequiredValidator(fl validator.FieldLevel) bool {
	if isParentExist(fl) && fl.Field().String() == "" {
		return false
	}
	return true
}

func RegisterDocumentValidator() *validator.Validate {
	v := validator.New()
	v.RegisterValidation("optionalRequired", optionalRequiredValidator)
	return v
}

func RegisterVcMetaValidator() *validator.Validate {
	v := validator.New()
	v.RegisterValidation("optionalRequired", optionalRequiredValidator)
	return v
}

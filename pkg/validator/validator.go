package validator

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	v *validator.Validate
}

func NewValidator() (*Validator, error) {
	v := validator.New()
	if err := v.RegisterValidation("amount", amountValidate); err != nil {
		return nil, err
	}
	return &Validator{v: v}, nil
}

func (v *Validator) Validate(i interface{}) error {
	if err := v.v.Struct(i); err != nil {
		return validateError(err.(validator.ValidationErrors)[0])
	}
	return nil
}

func validateError(err validator.FieldError) error {
	switch err.Tag() {
	case "amount":
		return errors.New("field amount is incorrect. Amount must be > 0")
	default:
		return fmt.Errorf("field %s is required and invalid", err.Field())
	}
}

func amountValidate(fl validator.FieldLevel) bool {
	return fl.Field().Float() > 0
}

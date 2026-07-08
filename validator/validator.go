package validator

import (
	"github.com/go-playground/validator/v10"
)

// CustomValidator wraps go-playground/validator for use with Echo.
type CustomValidator struct {
	validate *validator.Validate
}

func New() *CustomValidator {
	return &CustomValidator{validate: validator.New()}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validate.Struct(i)
}

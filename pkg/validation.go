package pkg

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// Validate implements validation for a Competition.
func (c Competition) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&a.Name, validation.Required),
	)
}

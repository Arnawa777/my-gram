package models

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// GORMModel represents the model for an templates
type GORMModel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at",omitempty`
	UpdatedAt time.Time `json:"updated_at",omitempty`
}

// This Below Code is not best practice to write it on Models

// Set json message when error
type validationError struct {
	Message string `json:"message"`
}

// to check validation when Create or Update
func GetValidationErrors(data interface{}) []validationError {
	errs := []validationError{}
	if _, err := govalidator.ValidateStruct(data); err != nil {
		for _, e := range err.(govalidator.Errors).Errors() {
			errs = append(errs, validationError{Message: e.Error()})
		}
	}
	return errs
}

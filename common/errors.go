package common

import (
	"encoding/json"
	"errors"
)

// FieldValidation struct
type FieldValidation struct {
	Field  string
	Errors []error
}

// ErrorValidation error
type ErrorValidation struct {
	Fields []FieldValidation
}

// AddError to validation
func (e *ErrorValidation) AddError(field string, err error) {
	var ex *FieldValidation

	// Find existing field
	var modfields []FieldValidation
	for _, f := range e.Fields {
		if f.Field == field {
			ex = &f
		} else {
			modfields = append(modfields, f)
		}
	}

	if ex == nil {
		modfields = append(e.Fields, FieldValidation{Field: field, Errors: []error{err}})
	} else {
		ex.Errors = append(ex.Errors, err)
		modfields = append(modfields, *ex)
	}

	e.Fields = modfields
}

// Get error
func (e *ErrorValidation) Get() error {
	if len(e.Fields) == 0 {
		return errors.New("Validation error")
	}

	// Map error
	type errstr struct {
		Field string   `json:"field"`
		Err   []string `json:"errors"`
	}
	msg := []errstr{}

	for _, field := range e.Fields {
		errfield := []string{}
		for _, err := range field.Errors {
			errfield = append(errfield, err.Error())
		}
		msg = append(msg, errstr{
			Field: field.Field,
			Err:   errfield,
		})
	}

	errs, _ := json.Marshal(msg)
	return errors.New(string(errs))

}

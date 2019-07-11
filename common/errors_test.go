package common

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationError(t *testing.T) {
	val := ErrorValidation{}

	val.AddError("email", errors.New("Required"))
	assert.Equal(t, `[{"field":"email","errors":["Required"]}]`, val.Get().Error())

	val.AddError("email", errors.New("Not Email"))

	assert.Equal(t, `[{"field":"email","errors":["Required","Not Email"]}]`, val.Get().Error())

	val.AddError("password", errors.New("Cannot blank"))

	assert.Equal(t, `[{"field":"email","errors":["Required","Not Email"]},{"field":"password","errors":["Cannot blank"]}]`, val.Get().Error())
}

func TestValidationEmpty(t *testing.T) {
	val := ErrorValidation{}

	assert.Equal(t, "Validation error", val.Get().Error())
}

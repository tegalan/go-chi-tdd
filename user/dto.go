package user

import (
	"errors"
	"go-chi/common"
	"net/http"
)

// SignUpRequest ...
type SignUpRequest struct {
	*User
}

// Bind to user
func (b *SignUpRequest) Bind(r *http.Request) error {
	v := common.ErrorValidation{}

	if b.User == nil {
		return errors.New("Empty request")
	}

	if b.User.Email == "" {
		v.AddError("email", errors.New("Email cannot empty"))
	}

	if b.User.Name == "" {
		v.AddError("name", errors.New("Name cannot empty"))
	}

	if b.User.Password == "" {
		v.AddError("password", errors.New("Password cannot empty"))
	}

	if len(v.Fields) > 0 {
		return v.Get()
	}

	return nil
}

// LogiInRequest struct
type LogiInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Bind login request
func (b *LogiInRequest) Bind(r *http.Request) error {
	if b.Email == "" {
		return errors.New("Email cannot empty")
	}

	if b.Password == "" {
		return errors.New("Password cannot empty")
	}

	return nil
}

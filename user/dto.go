package user

import (
	"errors"
	"net/http"
)

// SignUpRequest ...
type SignUpRequest struct {
	*User
}

// Bind to user
func (b *SignUpRequest) Bind(r *http.Request) error {
	if b.User == nil {
		return errors.New("Empty request")
	}

	if b.User.Email == "" {
		return errors.New("Email cannot empty")
	}

	if b.User.Name == "" {
		return errors.New("Name cannot empty")
	}

	if b.User.Password == "" {
		return errors.New("Password cannot empty")
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

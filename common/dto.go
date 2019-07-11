package common

import (
	"net/http"

	"github.com/go-chi/render"
)

// ErrorResponse struct
type ErrorResponse struct {
	Error   error  `json:"-"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Render implement render func
func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.Code)
	return nil
}

// ErrorNotFound error
func ErrorNotFound(err error) render.Renderer {
	return &ErrorResponse{
		Error:   err,
		Code:    http.StatusNotFound,
		Message: err.Error(),
	}
}

// ErrorUnprocessable error
func ErrorUnprocessable(err error) render.Renderer {
	return &ErrorResponse{
		Error:   err,
		Code:    http.StatusUnprocessableEntity,
		Message: err.Error(),
	}
}

// ErrorUnauthorized error
func ErrorUnauthorized(err error) render.Renderer {
	return &ErrorResponse{
		Error:   err,
		Code:    http.StatusUnauthorized,
		Message: err.Error(),
	}
}

// ErrorRender common error response
func ErrorRender(err error, code int) render.Renderer {
	return &ErrorResponse{
		Error:   err,
		Code:    code,
		Message: err.Error(),
	}
}

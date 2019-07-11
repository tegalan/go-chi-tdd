package user

import (
	"errors"
	"go-chi/common"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// Handler for user
type Handler struct {
	store Store
}

// Routes to mounted on chi
func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/", h.SignUp)

	return r
}

// SignUp handler
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	data := &SignUpRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, common.ErrorUnprocessable(err))
		return
	}

	user := data.User
	if err := h.store.Create(user); err != nil {
		render.Render(w, r, common.ErrorRender(errors.New("Unable to create user"), http.StatusBadRequest))
		return
	}

	render.JSON(w, r, user)
}

// Login handler
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	data := &LogiInRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, common.ErrorUnprocessable(err))
		return
	}

	u, e := h.store.FindByLogin(data.Email, data.Password)
	if e != nil {
		render.Render(w, r, common.ErrorUnauthorized(e))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, u)
}

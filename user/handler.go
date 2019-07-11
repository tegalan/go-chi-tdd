package user

import (
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
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{"message": "Error bad request"})
		return
	}

	user := data.User
	if err := h.store.Create(user); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{"message": "Unable to create user"})
		return
	}

	render.JSON(w, r, user)
}

// Login handler
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	data := &LogiInRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Status(r, http.StatusUnprocessableEntity)
		render.JSON(w, r, render.M{"message": err.Error()})
		return
	}

	u, e := h.store.FindByLogin(data.Email, data.Password)
	if e != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, render.M{"message": e.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, u)
}

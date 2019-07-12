package user

import (
	"errors"
	"go-chi/common"
	"go-chi/router"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// Handler for user
type Handler struct {
	store Store
}

// Register ...
func Register(s Store) Handler {
	h := Handler{
		store: s,
	}

	r := router.GetRouter()
	r.Mount("/user/", h.Routes())
	return h
}

// Routes to mounted on chi
func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/signup", h.SignUp)
	r.Post("/login", h.Login)

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
	// Handle payload
	data := &LogiInRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, common.ErrorUnprocessable(err))
		return
	}

	// Find user record from store
	u, e := h.store.FindByLogin(data.Email, data.Password)
	if e != nil {
		render.Render(w, r, common.ErrorUnauthorized(e))
		return
	}

	// Generate Token
	var token string
	render.Render(w, r, NewLoginResponse(u, token))
}

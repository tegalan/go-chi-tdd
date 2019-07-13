package user

import (
	"errors"
	"go-chi/common"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// Handler for user
type Handler struct {
	store Store
}

// NewHandler ...
func NewHandler(s Store) Handler {
	h := Handler{
		store: s,
	}
	return h
}

// RegisterRouter to mounted on chi
func (h *Handler) RegisterRouter(router *chi.Mux) http.Handler {
	r := chi.NewRouter()

	r.Post("/signup", h.SignUp)
	r.Post("/login", h.Login)

	router.Mount("/user", r)
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
		log.Println(err)
		render.Render(w, r, common.ErrorRender(errors.New("Unable to create user"), http.StatusBadRequest))
		return
	}

	render.Render(w, r, NewSignupRequest(user))
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
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Minute * 30).Unix(),
	})

	// TODO: create secure secret
	token, err := t.SignedString([]byte("rahasia"))

	if err != nil {
		render.Render(w, r, common.ErrorRender(errors.New("Failed when generate token"), http.StatusInternalServerError))
	}

	render.Render(w, r, NewLoginResponse(u, token))
}

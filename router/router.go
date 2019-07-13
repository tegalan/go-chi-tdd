package router

import (
	"errors"
	"fmt"
	"go-chi/common"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// GetRouter get router
func GetRouter() *chi.Mux {

	route := chi.NewMux()

	route.Use(middleware.Logger)
	route.Use(Recover)

	route.MethodNotAllowed(MethodNotAllowedHandler)
	route.NotFound(NotFoundHandler)

	route.Get("/test/500", func(w http.ResponseWriter, r *http.Request) {
		panic("Oops, recover me!")
	})

	return route
}

// Recover custom recover middleware
func Recover(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {

				logEntry := middleware.GetLogEntry(r)
				if logEntry != nil {
					logEntry.Panic(rvr, debug.Stack())
				} else {
					fmt.Fprintf(os.Stderr, "Panic: %+v\n", rvr)
					debug.PrintStack()
				}

				render.Render(w, r, common.ErrorRender(errors.New("Internal server error"), http.StatusInternalServerError))
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// NotFoundHandler ...
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, common.ErrorRender(errors.New("Not Found"), http.StatusNotFound))
}

// MethodNotAllowedHandler ...
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, common.ErrorRender(errors.New("Method not allowed"), http.StatusMethodNotAllowed))
}

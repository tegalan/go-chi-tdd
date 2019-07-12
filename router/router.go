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

var route *chi.Mux

func init() {
	route = chi.NewMux()
	route.Use(middleware.Logger)
	//route.Use(middleware.Recoverer)
	route.Use(Recover)
	route.MethodNotAllowed(MethodNotAllowedHandler)
	route.NotFound(NotFoundHandler)

	route.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})

	route.Get("/500", func(w http.ResponseWriter, r *http.Request) {
		panic("Oops, recover me!")
	})
}

// GetRouter get router
func GetRouter() *chi.Mux {
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

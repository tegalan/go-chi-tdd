package main

import (
	"database/sql"
	"errors"
	"go-chi/common"
	"go-chi/user"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	_ "github.com/lib/pq"
)

func main() {

	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/go_chi?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		render.Render(w, r, common.ErrorRender(errors.New("Method not allowed"), http.StatusMethodNotAllowed))
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		render.Render(w, r, common.ErrorRender(errors.New("Not Found"), http.StatusNotFound))
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})

	r.Get("/500", func(w http.ResponseWriter, r *http.Request) {
		panic("Oops, recover me!")
	})

	// User Resource
	user := user.NewHandler(&user.Model{Db: db})
	r.Mount("/user", user.Routes())

	http.ListenAndServe(":8000", r)
}

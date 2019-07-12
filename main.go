package main

import (
	"database/sql"
	"go-chi/app/user"
	"go-chi/router"
	"log"
	"net/http"

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

	// User Resource
	user.Register(&user.DatabaseStore{DB: db})
	//r.Mount("/user", user.Routes())

	r := router.GetRouter()
	http.ListenAndServe(":8000", r)
}

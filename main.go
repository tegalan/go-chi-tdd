package main

import (
	"database/sql"
	"go-chi/app/user"
	"go-chi/database"
	"go-chi/database/migrations"
	"go-chi/router"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {

	db, err := database.GetDB()
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Handle Migration
	if handleMigrateArgs(db) {
		return
	}

	// Register User Resource
	user.Register(&user.DatabaseStore{DB: db})
	//r.Mount("/user", user.Routes())

	r := router.GetRouter()
	http.ListenAndServe(":8000", r)
}

func handleMigrateArgs(db *sql.DB) bool {

	args := os.Args[1:]

	if len(args) == 2 && args[0] == "migrate" && args[1] == "down" {
		migrations.MigrateDown(db)

		return true
	}

	if len(args) == 1 && args[0] == "migrate" {
		migrations.MigrateUp(db)
		return true
	}

	return false
}

package database

import (
	"database/sql"
)

var db *sql.DB

// GetDB get connection
func GetDB() (*sql.DB, error) {
	if db != nil {
		return db, nil
	}

	var err error
	db, err = sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/go_chi?sslmode=disable")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil

}

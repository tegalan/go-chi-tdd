package database

import (
	"database/sql"
	"go-chi/conf"
)

var db *sql.DB

// GetDB get connection
func GetDB() (*sql.DB, error) {
	if db != nil {
		return db, nil
	}

	// Get config for default use
	conf := conf.GetConfig(false)

	var err error
	db, err = sql.Open("postgres", conf.DBUrl)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil

}

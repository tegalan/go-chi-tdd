package migrations

import (
	"database/sql"
)

func init() {
	AddMigration(1, up0001, down0001)
}

func up0001(tx *sql.Tx) error {
	// Create table schema
	_, err := tx.Exec(`CREATE TABLE users (id SERIAL PRIMARY KEY, name VARCHAR NULL, email VARCHAR NOT NULL, password VARCHAR NOT NULL)`)

	if err != nil {
		//tx.Rollback()
		return err
	}

	// Also can do seed database record
	_, ierr := tx.Exec("INSERT INTO users (name, email, password) VALUES($1, $2, $3)", "Moana", "moana@motunui.is", "$2y$12$jdaF7iJqKt2MGDFOrJGtaOvz.sOtE4E/IdMiflzztp.a5za4u6KhO")
	if ierr != nil {
		return err
	}

	//tx.Commit()
	return nil
}

func down0001(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE users")
	if err != nil {
		//tx.Rollback()
		return err
	}

	//tx.Commit()
	return nil
}

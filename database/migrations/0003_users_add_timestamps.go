package migrations

import "database/sql"

func init() {
	AddMigration(3, up0003, down0003)
}

func up0003(tx *sql.Tx) error {
	if _, err := tx.Exec("ALTER TABLE users ADD created_at TIMESTAMP NOT NULL DEFAULT NOW(), ADD updated_at TIMESTAMP NULL"); err != nil {
		return err
	}

	return nil
}

func down0003(tx *sql.Tx) error {
	if _, err := tx.Exec("ALTER TABLE users DROP COLUMN created_at, DROP COLUMN updated_at"); err != nil {
		return err
	}

	return nil
}

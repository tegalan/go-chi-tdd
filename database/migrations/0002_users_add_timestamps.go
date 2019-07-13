package migrations

import "database/sql"

func init() {
	AddMigration(2, up0002, down0002)
}

func up0002(tx *sql.Tx) error {
	if _, err := tx.Exec("ALTER TABLE users ADD created_at TIMESTAMP NOT NULL DEFAULT NOW(), ADD updated_at TIMESTAMP NULL"); err != nil {
		return err
	}

	return nil
}

func down0002(tx *sql.Tx) error {
	if _, err := tx.Exec("ALTER TABLE users DROP COLUMN created_at, DROP COLUMN updated_at"); err != nil {
		return err
	}

	return nil
}

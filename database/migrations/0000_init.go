package migrations

import (
	"database/sql"
)

func up(tx *sql.Tx) error {
	return nil
}

func down(tx *sql.Tx) error {
	return nil
}

func init() {
	AddMigration(0, up, down)
}

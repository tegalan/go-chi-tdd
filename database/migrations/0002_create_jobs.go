package migrations

import "database/sql"

func init() {
	AddMigration(2, up0002, down0002)
}

func up0002(tx *sql.Tx) error {
	return nil
}

func down0002(tx *sql.Tx) error {
	return nil
}

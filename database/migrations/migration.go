package migrations

import (
	"database/sql"
	"log"
	"sort"

	"github.com/lib/pq"
)

var migrates []Migration

// Migration struct
type Migration struct {
	Version int
	Up      func(tx *sql.Tx) error
	Down    func(tx *sql.Tx) error
}

// History struct
type History struct {
	version int
	step    int
}

// AddMigration function to array
func AddMigration(version int, up func(tx *sql.Tx) error, down func(tx *sql.Tx) error) {
	migrates = append(migrates, Migration{
		Version: version,
		Up:      up,
		Down:    down,
	})
}

// InitMigration Check migrations table
func InitMigration(db *sql.DB) {
	row := 0
	err := db.QueryRow("SELECT COUNT(*) FROM migrations").Scan(&row)
	if err != nil {
		if e, ok := err.(*pq.Error); ok {
			if e.Code.Name() == "undefined_table" {
				createMigrationsTable(db)
				return
			}
		}
	}
}

// createMigrationsTable ...
func createMigrationsTable(db *sql.DB) {
	log.Println("0000 - Initializing migrations table")
	_, err := db.Exec("CREATE TABLE migrations (version SERIAL PRIMARY KEY, step INTEGER NOT NULL)")
	if err != nil {
		log.Println("0000 - Unable to create migrations table", err)
		return
	}
}

func getHistory(db *sql.DB) []History {
	var histories []History

	row, err := db.Query("SELECT version, step FROM migrations")
	if err != nil {
		log.Println("Unable to read log history")
	}

	defer row.Close()
	for row.Next() {
		var history History
		err := row.Scan(&history.version, &history.step)
		if err != nil {
			log.Fatal(err)
		}

		histories = append(histories, history)
	}

	return histories
}

func isMigrated(histories []History, version int) bool {
	for _, h := range histories {
		if h.version == version {
			return true
		}
	}
	return false
}

func isForRollback(histories []History, step int, ver int) bool {

	for _, h := range histories {
		if h.version == ver && h.step == step {
			return true
		}
	}

	return false
}

func getMaxStep(db *sql.DB) int {
	step := 0

	err := db.QueryRow("SELECT COALESCE(MAX(step), 0) FROM migrations").Scan(&step)
	if err != nil {
		log.Fatal(err)
	}

	return step
}

// MigrateUp ...
func MigrateUp(db *sql.DB) error {
	// Check if migrations table is initalized
	InitMigration(db)

	// Sort by version asc
	sort.Slice(migrates, func(i, j int) bool {
		return migrates[i].Version < migrates[j].Version
	})

	// Migration history
	h := getHistory(db)

	// Get max step
	step := getMaxStep(db) + 1

	counter := 0
	for _, mig := range migrates {
		// Check if already migrated
		if !isMigrated(h, mig.Version) {
			tx, err := db.Begin()
			if err != nil {
				return err
			}
			log.Printf("%04d - Migrating \n", mig.Version)

			if err := mig.Up(tx); err != nil {
				tx.Rollback()
				log.Printf("%04d - Error migrate: %s\n", mig.Version, err.Error())
				return err
			}

			// Save migration history
			_, e := tx.Exec("INSERT INTO migrations (version, step) VALUES($1, $2)", mig.Version, step)
			if e != nil {
				tx.Rollback()
				log.Printf("%04d - Error migrate: %s\n", mig.Version, e.Error())
				return e
			}

			tx.Commit()
			counter++
		}

	}

	if counter > 0 {
		log.Println("Database migration success...")
	} else {
		log.Println("Nothing to migrate!")
	}
	return nil
}

// MigrateDown ...
func MigrateDown(db *sql.DB) error {

	// Sort by version desc
	sort.Slice(migrates, func(i, j int) bool {
		return migrates[i].Version > migrates[j].Version
	})

	// Migration history
	h := getHistory(db)

	// Get max step
	step := getMaxStep(db)

	counter := 0
	for _, mig := range migrates {
		// Check if already migrated & 1 step for rollback
		if isMigrated(h, mig.Version) && isForRollback(h, step, mig.Version) {
			log.Printf("%04d - Rolling back \n", mig.Version)

			tx, err := db.Begin()

			if err != nil {
				return err
			}

			if err := mig.Down(tx); err != nil {
				tx.Rollback()
				log.Printf("%04d - Error rollback: %s\n", mig.Version, err.Error())
				return err
			}

			// Delete Migration history
			_, e := tx.Exec("DELETE FROM migrations WHERE version = $1", mig.Version)

			if e != nil {
				tx.Rollback()
				log.Printf("%04d - Error rollback: %s\n", mig.Version, e.Error())
				return e
			}
			tx.Commit()
			counter++
		}
	}

	if counter > 0 {
		log.Println("Database rollback success...")
	} else {
		log.Println("Nothing to rollback!")
	}
	return nil
}

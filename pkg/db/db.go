package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `CREATE TABLE scheduler (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				date CHAR(8) NOT NULL DEFAULT "",
				title TEXT NOT NULL,
				comment TEXT,
				repeat VARCHAR(128)
				);`

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	var db *sql.DB
	var install bool

	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}
	defer db.Close()

	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return fmt.Errorf("failed create a database: %w", err)
		}
	}

	return nil
}

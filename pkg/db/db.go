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

var DB *sql.DB

func Init(dbFile string) error {
	var install bool
	var err error
	if _, err = os.Stat(dbFile); os.IsNotExist(err) {
		install = true
	} else if err != nil {
		return fmt.Errorf("stat %q failed: %w", dbFile, err)
	}

	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	if install {
		_, err = DB.Exec(schema)
		if err != nil {
			return fmt.Errorf("failed create a database: %w", err)
		}
	}

	return nil
}

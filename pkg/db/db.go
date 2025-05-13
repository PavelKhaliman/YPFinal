package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const (
	schema = `CREATE TABLE scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8) NOT NULL DEFAULT "",
			title VARCHAR,
			comment TEXT,
			repeat VARCHAR
	);
	CREATE INDEX idx_scheduler_date ON scheduler (date);`
)

func Init(dbFile string) error {

	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {

		install = true
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка при открытии базы данных: %w", err)
	}

	if install {
		_, err := db.Exec(schema)
		if err != nil {
			db.Close()
			return fmt.Errorf("ошибка при создании схемы базы данных: %w", err)
		}
	}

	return nil
}

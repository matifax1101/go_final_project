package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var database *sql.DB

const createTable = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
)`

const createIndex = `CREATE INDEX idx_scheduler_date ON scheduler(date)`

func Init(dbFile string) error {
	if database != nil {
		return nil
	}

	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)
	if err != nil && !install {
		return err
	}

	database, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if err = database.Ping(); err != nil {
		database.Close()
		database = nil
		return err
	}

	if install {
		if _, err = database.Exec(createTable); err != nil {
			database.Close()
			database = nil
			return fmt.Errorf("create table: %w", err)
		}
		if _, err = database.Exec(createIndex); err != nil {
			database.Close()
			database = nil
			return fmt.Errorf("create index: %w", err)
		}
	}

	return nil
}

func Close() {
	if database != nil {
		_ = database.Close()
		database = nil
	}
}

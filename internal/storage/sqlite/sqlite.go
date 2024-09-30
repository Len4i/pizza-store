package sqlite

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func New(storagePath string) (*sql.DB, error) {

	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS orders (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		size TEXT NOT NULL,
		amount NUMBER NOT NULL,
		pizza_type TEXT NOT NULL);
	`)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, err
	}

	return db, nil
}

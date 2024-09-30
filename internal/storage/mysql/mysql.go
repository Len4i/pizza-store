package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func New(host string, port int, user string, pass string, dbName string) (*sql.DB, error) {

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?tls=skip-verify", user, pass, host, port, dbName)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("unable to open: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("unable to ping: %w", err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS orders (
		id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
		size TEXT NOT NULL,
		amount INTEGER NOT NULL,
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

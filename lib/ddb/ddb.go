package ddb

import (
	"database/sql"
)

func open(driver, connect string) (*sql.DB, error) {
	db, err := sql.Open(driver, connect)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	// db.Close()
	return db, nil
}

//Open sql
func Open(driver, connect string) (*sql.DB, error) {
	return open(driver, connect)
}

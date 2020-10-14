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
	db.SetConnMaxLifetime(0) //这里设置一下超时
	return db, nil
}

//Open sql
func Open(driver, connect string) (*sql.DB, error) {
	return open(driver, connect)
}

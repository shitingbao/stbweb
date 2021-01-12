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

//dropLikeTable 删除数据库内，表名称包含‘like’内容的表
//实际先使用第一句查出所有表，拼接成drop语句，第二步执行
//未使用事务，待定
func dropLikeTable(like string, db *sql.DB) error {
	rows, err := db.Query(`
	Select CONCAT( 'drop table ', table_name, ';' ) 
		FROM information_schema.tables 
	Where table_schema='stbweb' and table_name LIKE ?`, "%"+like+"%")
	if err != nil {
		return err
	}
	dropSQL := ""
	for rows.Next() {
		if err := rows.Scan(&dropSQL); err != nil {
			return err
		}
		if dropSQL == "" {
			continue
		}

		stmt, err := db.Prepare(dropSQL)
		if err != nil {
			return err
		}
		if _, err := stmt.Exec(); err != nil {
			return err
		}
	}
	return nil
}

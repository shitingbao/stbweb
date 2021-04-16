package server

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	db, err := openMysql()
	if err != nil {
		panic(err)
	}
	DB = db
}

//"root", "12345678", "127.0.0.1", "test", "3306"
func openMysql() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open("root:12345678@(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})
	// db, err := gorm.Open(mysql.Open("test:yZx0nqxbwQqv@(sh-cdb-mr97xubq.sql.tencentcdb.com:61637)/redkol_test?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		log.Println("orm:", err)
		return db, err
	}
	return db, nil
}

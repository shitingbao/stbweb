package ddb

import (
	"github.com/jinzhu/gorm"
)

//SQLOpen sqlopen
func SQLOpen(isLog bool) *gorm.DB {
	var err error
	db, err := gorm.Open("mysql", "root:123456@(127.0.0.1:3306)/ep?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	db.LogMode(isLog)
	return db
}

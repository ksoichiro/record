package db

import (
	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func Init() {
	var err error
	db, err = gorm.Open("mysql", "app:app@tcp(127.0.0.1:3306)/test?parseTime=true&loc=Asia%2FTokyo")
	if err != nil {
		panic(err.Error())
	}
}

func GetDB() *gorm.DB {
	return db
}

func CloseDB() {
	db.Close()
}

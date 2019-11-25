package db

import (
	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/ksoichiro/record/config"
)

var db *gorm.DB

// Init initializes the database handles.
func Init() {
	var err error
	db, err = gorm.Open(
		config.GetConfig().GetString("database.dialect"),
		config.GetConfig().GetString("database.url"))
	if err != nil {
		panic(err.Error())
	}
}

// SetDB sets the alternative gorm.DB instance mainly for the tests.
func SetDB(alternative *gorm.DB) {
	db = alternative
}

// GetDB returns the database handle.
func GetDB() *gorm.DB {
	return db
}

// CloseDB closes the database handle.
func CloseDB() {
	db.Close()
}

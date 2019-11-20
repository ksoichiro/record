package db

import (
	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

// Init initializes the database handles.
func Init() {
	var err error
	db, err = gorm.Open("mysql", "app:app@tcp(127.0.0.1:3306)/test?parseTime=true&loc=Asia%2FTokyo")
	if err != nil {
		panic(err.Error())
	}
}

// InitForTest initializes the in-memory database for the tests.
func InitForTest() {
	testDB, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	SetDB(testDB)
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

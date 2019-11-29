package models

import (
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ksoichiro/record/api/config"
	"github.com/ksoichiro/record/api/db"
	"github.com/ksoichiro/record/api/forms"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestNewRecord(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&User{}, &Task{}, &Record{})
	db.Create(&User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	db.Create(&Task{ID: 200, UserID: 100, Title: "task1", Description: "task description", Type: 0, Amount: 0, CreatedAt: time.Now()})
	form := forms.RecordCreateForm{TaskID: new(int)}
	*form.TaskID = 200
	userID := 100
	record, err := NewRecord(&form, userID, mustParse("2006-01-02", "2019-11-20"))
	assert.Nil(t, err)
	assert.Equal(t, 100, record.User.ID)
	assert.Equal(t, 200, record.Task.ID)
	assert.Equal(t, "2019-11-20", record.TargetDate.Format("2006-01-02"))
	var records []Record
	var count int
	db.Find(&records).Count(&count)
	assert.Equal(t, 1, count)
	assert.Equal(t, 1, records[0].ID)
	assert.Equal(t, 100, records[0].UserID)
	assert.Equal(t, 200, records[0].TaskID)
	assert.Equal(t, "2019-11-20", records[0].TargetDate.Format("2006-01-02"))
}

func TestNewRecordFailsDueToTaskNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&User{}, &Task{}, &Record{})
	db.Create(&User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	db.Create(&Task{ID: 200, UserID: 101, Title: "task1", Description: "task description", Type: 0, Amount: 0, CreatedAt: time.Now()})
	form := forms.RecordCreateForm{TaskID: new(int)}
	*form.TaskID = 200
	userID := 100
	_, err := NewRecord(&form, userID, mustParse("2006-01-02", "2019-11-20"))
	assert.Equal(t, "record not found", err.Error())
}

func TestNewRecordFailsWhenAlreadyCreated(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&User{}, &Task{}, &Record{})
	db.Create(&User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	db.Create(&Task{ID: 200, UserID: 100, Title: "task1", Description: "task description", Type: 0, Amount: 0, CreatedAt: time.Now()})
	form := forms.RecordCreateForm{TaskID: new(int)}
	*form.TaskID = 200
	userID := 100
	_, err := NewRecord(&form, userID, mustParse("2006-01-02", "2019-11-20"))
	assert.Nil(t, err)
	_, err = NewRecord(&form, userID, mustParse("2006-01-02", "2019-11-20"))
	assert.Equal(t, "already created", err.Error())
}

func TestListRecords(t *testing.T) {
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&User{}, &Task{}, &Record{})
	db.Create(&User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	db.Create(&Task{ID: 200, UserID: 100, Title: "task1", Description: "task description", Type: 0, Amount: 0, CreatedAt: time.Now()})
	db.Create(&Task{ID: 201, UserID: 101, Title: "task2", Description: "task description", Type: 0, Amount: 0, CreatedAt: time.Now()})
	db.Create(&Task{ID: 202, UserID: 100, Title: "task3", Description: "task description", Type: 0, Amount: 0, CreatedAt: time.Now()})
	db.Create(&Record{ID: 300, UserID: 100, TargetDate: mustParse("2006-01-02", "2019-11-19"), TaskID: 200, CreatedAt: time.Now()})
	db.Create(&Record{ID: 301, UserID: 100, TargetDate: mustParse("2006-01-02", "2019-11-20"), TaskID: 200, CreatedAt: time.Now()})
	db.Create(&Record{ID: 302, UserID: 101, TargetDate: mustParse("2006-01-02", "2019-11-20"), TaskID: 201, CreatedAt: time.Now()})
	db.Create(&Record{ID: 303, UserID: 100, TargetDate: mustParse("2006-01-02", "2019-11-21"), TaskID: 200, CreatedAt: time.Now()})
	db.Create(&Record{ID: 304, UserID: 100, TargetDate: mustParse("2006-01-02", "2019-11-20"), TaskID: 202, CreatedAt: time.Now()})
	records := ListRecords(100, mustParse("2006-01-02", "2019-11-20"))
	assert.Equal(t, 2, len(records))
	assert.Equal(t, 301, records[0].ID)
	assert.Equal(t, 304, records[1].ID)
}

func TestFindRecord(t *testing.T) {
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&User{}, &Task{}, &Record{})
	db.Create(&User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	db.Create(&Task{ID: 200, UserID: 100, Title: "task1", Description: "task description", Type: 0, Amount: 0, CreatedAt: time.Now()})
	db.Create(&Task{ID: 201, UserID: 101, Title: "task2", Description: "task description", Type: 0, Amount: 0, CreatedAt: time.Now()})
	db.Create(&Task{ID: 202, UserID: 100, Title: "task3", Description: "task description", Type: 0, Amount: 0, CreatedAt: time.Now()})
	db.Create(&Record{ID: 300, UserID: 100, TargetDate: mustParse("2006-01-02", "2019-11-19"), TaskID: 200, CreatedAt: time.Now()})
	db.Create(&Record{ID: 301, UserID: 100, TargetDate: mustParse("2006-01-02", "2019-11-20"), TaskID: 200, CreatedAt: time.Now()})
	db.Create(&Record{ID: 302, UserID: 101, TargetDate: mustParse("2006-01-02", "2019-11-20"), TaskID: 201, CreatedAt: time.Now()})
	db.Create(&Record{ID: 303, UserID: 100, TargetDate: mustParse("2006-01-02", "2019-11-21"), TaskID: 200, CreatedAt: time.Now()})
	db.Create(&Record{ID: 304, UserID: 100, TargetDate: mustParse("2006-01-02", "2019-11-20"), TaskID: 202, Done: true, CreatedAt: time.Now()})
	record, err := FindRecord(304, 100)
	assert.Nil(t, err)
	assert.Equal(t, 304, record.ID)
	assert.Equal(t, true, record.Done)
}

func TestFindRecordErrorRecordNotFound(t *testing.T) {
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&User{}, &Task{}, &Record{})
	db.Create(&User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	db.Create(&Task{ID: 200, UserID: 100, Title: "task1", Description: "task description", Type: 0, Amount: 0, CreatedAt: time.Now()})
	db.Create(&Task{ID: 201, UserID: 101, Title: "task2", Description: "task description", Type: 0, Amount: 0, CreatedAt: time.Now()})
	db.Create(&Task{ID: 202, UserID: 100, Title: "task3", Description: "task description", Type: 0, Amount: 0, CreatedAt: time.Now()})
	record, err := FindRecord(304, 100)
	assert.Equal(t, "record not found", err.Error())
	assert.Nil(t, record)
}

func TestFindRecordErrorRecordExistsButNotOwnedByUser(t *testing.T) {
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&User{}, &Task{}, &Record{})
	db.Create(&User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	db.Create(&Task{ID: 200, UserID: 100, Title: "task1", Description: "task description", Type: 0, Amount: 0, CreatedAt: time.Now()})
	db.Create(&Task{ID: 201, UserID: 101, Title: "task2", Description: "task description", Type: 0, Amount: 0, CreatedAt: time.Now()})
	db.Create(&Task{ID: 202, UserID: 100, Title: "task3", Description: "task description", Type: 0, Amount: 0, CreatedAt: time.Now()})
	db.Create(&Record{ID: 300, UserID: 100, TargetDate: mustParse("2006-01-02", "2019-11-19"), TaskID: 200, CreatedAt: time.Now()})
	db.Create(&Record{ID: 301, UserID: 100, TargetDate: mustParse("2006-01-02", "2019-11-20"), TaskID: 200, CreatedAt: time.Now()})
	db.Create(&Record{ID: 302, UserID: 101, TargetDate: mustParse("2006-01-02", "2019-11-20"), TaskID: 201, CreatedAt: time.Now()})
	db.Create(&Record{ID: 303, UserID: 100, TargetDate: mustParse("2006-01-02", "2019-11-21"), TaskID: 200, CreatedAt: time.Now()})
	db.Create(&Record{ID: 304, UserID: 100, TargetDate: mustParse("2006-01-02", "2019-11-20"), TaskID: 202, Done: true, CreatedAt: time.Now()})
	record, err := FindRecord(302, 100)
	assert.Equal(t, "record not found", err.Error())
	assert.Nil(t, record)
}

func mustParse(layout string, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

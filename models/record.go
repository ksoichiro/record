package models

import (
	"fmt"
	"time"

	"github.com/ksoichiro/record/db"
	"github.com/ksoichiro/record/forms"
)

// Record represents daily record of a task.
type Record struct {
	ID         int       `json:"id" gorm:"primary_key;auto_increment"`
	UserID     int       `json:"-" gorm:"not null"`
	User       User      `json:"-"`
	TargetDate time.Time `json:"target_date" gorm:"not null"`
	TaskID     int       `json:"-" gorm:"not null"`
	Task       Task      `json:"-"`
	Done       bool      `json:"done" gorm:"not null;default 0"`
	Amount     int       `json:"amount"`
	CreatedAt  time.Time `json:"created_at" gorm:"not null"`
}

// NewRecord creates a new record object and persist it to the database.
func NewRecord(json *forms.RecordCreateForm, userID int, targetDateExpr string) (*Record, error) {
	targetDate, err := time.Parse("2006-01-02", targetDateExpr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date")
	}
	db := db.GetDB()
	tx := db.Begin()
	var user User
	tx.Where("id = ?", userID).Find(&user)
	var task Task
	var count int
	tx.Where("id = ? and user_id = ?", *json.TaskID, userID).First(&task).Count(&count)
	if count == 0 {
		tx.Rollback()
		return nil, fmt.Errorf("task not found")
	}
	tx.Model(&Record{}).Where("user_id = ? and target_date = ? and task_id = ?", userID, targetDate, json.TaskID).Count(&count)
	if 0 < count {
		tx.Rollback()
		return nil, fmt.Errorf("already created")
	}
	record := Record{
		User:       user,
		TargetDate: targetDate,
		Task:       task,
		CreatedAt:  time.Now(),
	}
	if json.Done == nil {
		record.Done = false
	} else {
		record.Done = *json.Done
	}
	if json.Amount == nil {
		record.Amount = 0
	} else {
		record.Amount = *json.Amount
	}
	tx.Create(&record)
	tx.Commit()
	return &record, nil
}

// ListRecords gets the records by user ID and the specified date.
func ListRecords(userID int, targetDate time.Time) []Record {
	db := db.GetDB()
	records := []Record{}
	db.Where("user_id = ? and target_date = ?", userID, targetDate).Find(&records)
	return records
}

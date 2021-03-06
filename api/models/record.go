package models

import (
	"fmt"
	"time"

	"github.com/ksoichiro/record/api/db"
	"github.com/ksoichiro/record/api/forms"
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
func NewRecord(json *forms.RecordCreateForm, userID int, targetDate time.Time) (record *Record, err error) {
	db := db.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			record = nil
			err = fmt.Errorf("failed to create record")
		}
	}()
	var user User
	if err = tx.Where("id = ?", userID).Find(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	var task Task
	var count int
	if err = tx.Where("id = ? and user_id = ?", *json.TaskID, userID).First(&task).Count(&count).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if count == 0 {
		tx.Rollback()
		return nil, fmt.Errorf("record not found")
	}
	if err = tx.Model(&Record{}).Where("user_id = ? and target_date = ? and task_id = ?", userID, targetDate, json.TaskID).Count(&count).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if 0 < count {
		tx.Rollback()
		return nil, fmt.Errorf("already created")
	}
	record = &Record{
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
	if err = tx.Create(&record).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return record, nil
}

// ListRecords gets the records by user ID and the specified date.
func ListRecords(userID int, targetDate time.Time) []Record {
	db := db.GetDB()
	records := []Record{}
	db.Where("user_id = ? and target_date = ?", userID, targetDate).Find(&records)
	return records
}

// FindRecord finds the record specified by ID and owned by the user.
func FindRecord(id int, userID int) (*Record, error) {
	db := db.GetDB()
	var record Record
	var count int
	db.Where("id = ? and user_id = ?", id, userID).First(&record).Count(&count)
	if count == 0 {
		return nil, fmt.Errorf("record not found")
	}
	return &record, nil
}

// Update updates the existing record.
func (r Record) Update(json *forms.RecordUpdateForm) error {
	db := db.GetDB()
	tx := db.Begin()
	var record Record
	var count int
	tx.Where("id = ? and user_id = ?", r.ID, r.UserID).First(&record).Count(&count)
	if count == 0 {
		tx.Rollback()
		return fmt.Errorf("record not found")
	}
	if json.Done != nil {
		record.Done = *json.Done
	}
	if json.Amount != nil {
		record.Amount = *json.Amount
	}
	tx.Save(&record)
	tx.Commit()
	return nil
}

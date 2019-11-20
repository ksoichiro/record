package models

import (
	"time"

	"github.com/ksoichiro/record/db"
	"github.com/ksoichiro/record/forms"
)

// Task represents a task that may be repeatedly executed by someone.
type Task struct {
	ID          int       `json:"id"`
	UserID      int       `json:"-" gorm:"not null"`
	User        User      `json:"-"`
	Title       string    `json:"title" gorm:"size:200;not null"`
	Description string    `json:"description"`
	Done        bool      `json:"done" gorm:"not null;default 0"`
	Type        int       `json:"type" gorm:"not null"`
	Amount      int       `json:"amount"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null"`
}

// NewTask creates a new task and persist it to the database.
func NewTask(json *forms.TaskCreateForm, userID int) (*Task, error) {
	db := db.GetDB()
	tx := db.Begin()
	var user User
	tx.Where("id = ?", userID).Find(&user)
	task := Task{
		User:        user,
		Title:       json.Title,
		Description: json.Description,
		CreatedAt:   time.Now(),
	}
	if json.Done == nil {
		task.Done = false
	} else {
		task.Done = *json.Done
	}
	if json.Type == nil {
		task.Type = 0
	} else {
		task.Type = *json.Type
	}
	if json.Amount == nil {
		task.Amount = 0
	} else {
		task.Amount = *json.Amount
	}
	tx.Create(&task)
	tx.Commit()
	return &task, nil
}

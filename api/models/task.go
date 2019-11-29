package models

import (
	"fmt"
	"time"

	"github.com/ksoichiro/record/api/db"
	"github.com/ksoichiro/record/api/forms"
)

// Task represents a task that may be repeatedly executed by someone.
type Task struct {
	ID          int       `json:"id"`
	UserID      int       `json:"-" gorm:"not null"`
	User        User      `json:"-"`
	Title       string    `json:"title" gorm:"size:200;not null"`
	Description string    `json:"description"`
	Type        int       `json:"type" gorm:"not null"`
	Amount      int       `json:"amount"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null"`
}

// NewTask creates a new task and persist it to the database.
func NewTask(json *forms.TaskCreateForm, userID int) (task *Task, err error) {
	db := db.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			task = nil
			err = fmt.Errorf("failed to create task")
		}
	}()
	var user User
	tx.Where("id = ?", userID).Find(&user)
	task = &Task{
		User:        user,
		Title:       json.Title,
		Description: json.Description,
		CreatedAt:   time.Now(),
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
	if err = tx.Create(&task).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return task, nil
}

// ListTasks gets the tasks by user ID.
func ListTasks(userID int) []Task {
	db := db.GetDB()
	tasks := []Task{}
	db.Where("user_id = ?", userID).Find(&tasks)
	return tasks
}

// FindTask finds the task specified by ID and owned by the user.
func FindTask(id int, userID int) (*Task, error) {
	db := db.GetDB()
	var task Task
	var count int
	db.Where("id = ? and user_id = ?", id, userID).First(&task).Count(&count)
	if count == 0 {
		return nil, fmt.Errorf("task not found")
	}
	return &task, nil
}

// Update updates the existing task.
func (t Task) Update(json *forms.TaskUpdateForm) error {
	db := db.GetDB()
	tx := db.Begin()
	var task Task
	var count int
	tx.Where("id = ? and user_id = ?", t.ID, t.UserID).First(&task).Count(&count)
	if count == 0 {
		tx.Rollback()
		return fmt.Errorf("task not found")
	}
	if json.Title != nil {
		task.Title = *json.Title
	}
	if json.Description != nil {
		task.Description = *json.Description
	}
	if json.Type != nil {
		task.Type = *json.Type
	}
	if json.Amount != nil {
		task.Amount = *json.Amount
	}
	tx.Save(&task)
	tx.Commit()
	return nil
}

package models

import "time"

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

package models

import "time"

type Record struct {
	ID         int       `json:"id" gorm:"primary_key;auto_increment"`
	UserID     int       `json:"-" gorm:"not null"`
	User       User      `json:"-"`
	TargetDate time.Time `json:"target_date" gorm:"not null"`
	TaskID     int       `json:"-" gorm:"not null"`
	Task       Task      `json:"-"`
	Amount     int       `json:"amount"`
	CreatedAt  time.Time `json:"created_at" gorm:"not null"`
}

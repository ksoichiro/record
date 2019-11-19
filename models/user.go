package models

import "time"

type User struct {
	ID        int       `json:"id" gorm:"primary_key;auto_increment"`
	Name      string    `json:"name" gorm:"size:100;not null;unique"`
	Password  string    `json:"password" gorm:"size:100;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
}

package models

import (
	"time"

	"github.com/ksoichiro/record/db"
	"github.com/ksoichiro/record/forms"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user for this system.
type User struct {
	ID        int       `json:"id" gorm:"primary_key;auto_increment"`
	Name      string    `json:"name" gorm:"size:100;not null;unique"`
	Password  string    `json:"password" gorm:"size:100;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
}

// NewUser creates a new user and persist it to the database.
func NewUser(json *forms.UserCreateForm) (*User, error) {
	name := json.Name
	hash, err := bcrypt.GenerateFromPassword([]byte(json.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	password := string(hash)
	db := db.GetDB()
	tx := db.Begin()
	user := User{Name: name, Password: password, CreatedAt: time.Now()}
	tx.Create(&user)
	tx.Commit()
	return &user, nil
}

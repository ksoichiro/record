package models

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ksoichiro/record/config"
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

// Login authenticates the user with the credentials specified by the form.
func Login(json *forms.UserLoginForm) (tokenString string, err error) {
	tokenString = ""
	db := db.GetDB()

	user := User{}
	db.Where("name = ?", json.Name).First(&user)
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(json.Password)); err != nil {
		err = fmt.Errorf("invalid name or password")
		return
	}

	signBytes, err := ioutil.ReadFile(config.GetConfig().GetString("auth.keys.private"))
	if err != nil {
		panic(err)
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		panic(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": "https://idp.example.com",
		"aud": "https://api.example.com",
		"sub": user.ID,
		"nbf": time.Now().UTC().Unix(),
		"exp": time.Now().UTC().Add(24 * time.Hour).Unix(),
	})
	tokenString, err = token.SignedString(signKey)
	return
}

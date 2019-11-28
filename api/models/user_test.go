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

func TestNewUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&User{})
	form := forms.UserCreateForm{
		Name:     "foo",
		Password: "test",
	}
	user, err := NewUser(&form)
	assert.Nil(t, err)
	assert.Equal(t, "foo", user.Name)
	var users []User
	var count int
	db.Find(&users).Count(&count)
	assert.Equal(t, 1, count)
	assert.Equal(t, 1, users[0].ID)
	assert.Equal(t, "foo", users[0].Name)
	assert.NotEmpty(t, users[0].Password)
	assert.GreaterOrEqual(t, time.Now().Unix(), users[0].CreatedAt.Unix())
}

func TestUserLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&User{})
	db.Create(&User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	form := forms.UserLoginForm{
		Name:     "foo",
		Password: "test",
	}
	token, err := Login(&form)
	assert.Nil(t, err)
	assert.NotEmpty(t, "token", token)
}

func TestUserLoginErrorInvalidName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&User{})
	db.Create(&User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	form := forms.UserLoginForm{
		Name:     "bar",
		Password: "test",
	}
	token, err := Login(&form)
	assert.Equal(t, "invalid name or password", err.Error())
	assert.Empty(t, token)
}

func TestUserLoginErrorInvalidPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&User{})
	db.Create(&User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	form := forms.UserLoginForm{
		Name:     "foo",
		Password: "wrong",
	}
	token, err := Login(&form)
	assert.Equal(t, "invalid name or password", err.Error())
	assert.Empty(t, token)
}

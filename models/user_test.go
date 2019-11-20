package models

import (
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ksoichiro/record/db"
	"github.com/ksoichiro/record/forms"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db.InitForTest()
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

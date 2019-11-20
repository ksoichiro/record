package models

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ksoichiro/record/db"
	"github.com/ksoichiro/record/forms"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestNewTask(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db.InitForTest()
	db := db.GetDB()
	db.AutoMigrate(&User{})
	form := forms.TaskCreateForm{
		Title: "important task",
	}
	userID := 1
	task, err := NewTask(&form, userID)
	assert.Nil(t, err)
	assert.Equal(t, "important task", task.Title)
}

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

func TestNewTask(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db.InitForTest()
	db := db.GetDB()
	db.AutoMigrate(&User{}, &Task{})
	db.Create(&User{ID: 1, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	form := forms.TaskCreateForm{
		Title: "important task",
	}
	userID := 1
	task, err := NewTask(&form, userID)
	assert.Nil(t, err)
	assert.Equal(t, "important task", task.Title)
	var tasks []Task
	var count int
	db.Where("user_id = 1").Find(&tasks).Count(&count)
	assert.Equal(t, 1, count)
	assert.Equal(t, 1, tasks[0].ID)
	assert.Equal(t, "important task", tasks[0].Title)
}

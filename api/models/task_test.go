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

func TestNewTask(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&User{}, &Task{})
	db.Create(&User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	form := forms.TaskCreateForm{
		Title: "important task",
	}
	userID := 100
	task, err := NewTask(&form, userID)
	assert.Nil(t, err)
	assert.Equal(t, "important task", task.Title)
	var tasks []Task
	var count int
	db.Where("user_id = 100").Find(&tasks).Count(&count)
	assert.Equal(t, 1, count)
	assert.Equal(t, 1, tasks[0].ID)
	assert.Equal(t, "important task", tasks[0].Title)
}

func TestListTasks(t *testing.T) {
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&User{}, &Task{})
	db.Create(&User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	db.Create(&Task{ID: 200, UserID: 100, Title: "task1", Description: "task description", Type: 0, Amount: 0, CreatedAt: time.Now()})
	db.Create(&Task{ID: 201, UserID: 101, Title: "task2", Description: "task description", Type: 0, Amount: 0, CreatedAt: time.Now()})
	db.Create(&Task{ID: 202, UserID: 100, Title: "task3", Description: "task description", Type: 0, Amount: 0, CreatedAt: time.Now()})
	records := ListTasks(100)
	assert.Equal(t, 2, len(records))
	assert.Equal(t, 200, records[0].ID)
	assert.Equal(t, 202, records[1].ID)
}

func TestFindTask(t *testing.T) {
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&User{}, &Task{})
	db.Create(&User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	db.Create(&Task{ID: 200, UserID: 100, Title: "task1", Description: "task description", Type: 0, Amount: 0, CreatedAt: time.Now()})
	db.Create(&Task{ID: 201, UserID: 101, Title: "task2", Description: "task description", Type: 0, Amount: 0, CreatedAt: time.Now()})
	db.Create(&Task{ID: 202, UserID: 100, Title: "task3", Description: "task description", Type: 0, Amount: 0, CreatedAt: time.Now()})
	task, err := FindTask(200, 100)
	assert.Nil(t, err)
	assert.Equal(t, 200, task.ID)
	assert.Equal(t, "task1", task.Title)
}

func TestFindTaskErrorTaskNotFound(t *testing.T) {
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&User{}, &Task{})
	db.Create(&User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	task, err := FindTask(200, 100)
	assert.Equal(t, "task not found", err.Error())
	assert.Nil(t, task)
}

func TestFindTaskErrorTaskExistsButNotOwnedByUser(t *testing.T) {
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&User{}, &Task{})
	db.Create(&User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	db.Create(&Task{ID: 200, UserID: 101, Title: "task1", Description: "task description", Type: 0, Amount: 0, CreatedAt: time.Now()})
	task, err := FindTask(200, 100)
	assert.Equal(t, "task not found", err.Error())
	assert.Nil(t, task)
}

func TestTaskUpdateSuccessfully(t *testing.T) {
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&User{}, &Task{})
	db.Create(&User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	db.Create(&Task{ID: 200, UserID: 100, Title: "task1", Description: "task description", Type: 0, Amount: 0, CreatedAt: time.Now()})
	var task Task
	db.Where("id = 200").First(&task)
	form := forms.TaskUpdateForm{ID: new(int), Title: new(string)}
	*form.ID = 200
	*form.Title = "modified"
	err := task.Update(&form)
	assert.Nil(t, err)
	db.Where("id = 200").First(&task)
	assert.Equal(t, "modified", task.Title)
}

func TestTaskUpdateErrorTaskNotFound(t *testing.T) {
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&User{}, &Task{})
	db.Create(&User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	task := Task{ID: 200, Title: "task1", Description: "task description", Type: 0, Amount: 0, CreatedAt: time.Now()}
	form := forms.TaskUpdateForm{ID: new(int), Title: new(string)}
	*form.ID = 200
	*form.Title = "modified"
	err := task.Update(&form)
	assert.Equal(t, "task not found", err.Error())
	var count int
	db.Where("id = 200").Count(&count)
	assert.Equal(t, 0, count)
}

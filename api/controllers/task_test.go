package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ksoichiro/record/api/config"
	"github.com/ksoichiro/record/api/db"
	"github.com/ksoichiro/record/api/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestTaskCreateSuccessfully(t *testing.T) {
	router := gin.Default()
	c := new(TaskController)
	gin.SetMode(gin.TestMode)
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&models.User{}, &models.Task{})
	db.Create(&models.User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	router.Use(func(c *gin.Context) {
		c.Set("user", 100)
	})
	router.POST("/create", c.Create)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/create", strings.NewReader(`{"title":"task1"}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"message":"created"}`, strings.TrimRight(w.Body.String(), "\n"))
}

func TestTaskCreateValidationError(t *testing.T) {
	router := gin.Default()
	c := new(TaskController)
	router.POST("/create", c.Create)
	gin.SetMode(gin.TestMode)
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&models.User{}, &models.Task{})
	db.Create(&models.User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/create", strings.NewReader(`{}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"error":"Key: 'TaskCreateForm.Title' Error:Field validation for 'Title' failed on the 'required' tag"}`, strings.TrimRight(w.Body.String(), "\n"))
}

func TestTaskCreateValidationErrorUserNotFound(t *testing.T) {
	router := gin.Default()
	c := new(TaskController)
	gin.SetMode(gin.TestMode)
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&models.User{}, &models.Task{})
	router.POST("/create", c.Create)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/create", strings.NewReader(`{"title":"task"}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, 403, w.Code, w.Body.String())
	assert.Equal(t, `{"error":"user not found"}`, strings.TrimRight(w.Body.String(), "\n"))
}

func TestTaskUpdateSuccessfully(t *testing.T) {
	router := gin.Default()
	c := new(TaskController)
	gin.SetMode(gin.TestMode)
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&models.User{}, &models.Task{})
	db.Create(&models.User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	db.Create(&models.Task{ID: 200, UserID: 100, Title: "task1", Description: "task description", Type: 0, Amount: 0, CreatedAt: time.Now()})
	router.Use(func(c *gin.Context) {
		c.Set("user", 100)
	})
	router.POST("/update", c.Update)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/update", strings.NewReader(`{"id":200,"title":"modified"}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"message":"updated"}`, strings.TrimRight(w.Body.String(), "\n"))
}

func TestTaskUpdateValidationError(t *testing.T) {
	router := gin.Default()
	c := new(TaskController)
	router.POST("/update", c.Update)
	gin.SetMode(gin.TestMode)
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&models.User{}, &models.Task{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/update", strings.NewReader(`{}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"error":"Key: 'TaskUpdateForm.ID' Error:Field validation for 'ID' failed on the 'exists' tag"}`, strings.TrimRight(w.Body.String(), "\n"))
}

func TestTaskUpdateValidationErrorUserNotFound(t *testing.T) {
	router := gin.Default()
	c := new(TaskController)
	gin.SetMode(gin.TestMode)
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&models.User{}, &models.Task{})
	router.POST("/update", c.Update)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/update", strings.NewReader(`{"id":200,"title":"modified"}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, 403, w.Code, w.Body.String())
	assert.Equal(t, `{"error":"user not found"}`, strings.TrimRight(w.Body.String(), "\n"))
}

func TestTaskUpdateErrorTaskNotFound(t *testing.T) {
	router := gin.Default()
	c := new(TaskController)
	gin.SetMode(gin.TestMode)
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&models.User{}, &models.Task{})
	db.Create(&models.User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	router.Use(func(c *gin.Context) {
		c.Set("user", 100)
	})
	router.POST("/update", c.Update)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/update", strings.NewReader(`{"id":200,"title":"modified"}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, 500, w.Code, w.Body.String())
	assert.Equal(t, `{"error":"task not found"}`, strings.TrimRight(w.Body.String(), "\n"))
}

func TestTaskUpdateErrorTaskExistsButNotOwnedByUser(t *testing.T) {
	router := gin.Default()
	c := new(TaskController)
	gin.SetMode(gin.TestMode)
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&models.User{}, &models.Task{})
	db.Create(&models.User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	db.Create(&models.Task{ID: 200, UserID: 101, Title: "task1", Description: "task description", Type: 0, Amount: 0, CreatedAt: time.Now()})
	router.Use(func(c *gin.Context) {
		c.Set("user", 100)
	})
	router.POST("/update", c.Update)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/update", strings.NewReader(`{"id":200,"title":"modified"}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, 500, w.Code, w.Body.String())
	assert.Equal(t, `{"error":"task not found"}`, strings.TrimRight(w.Body.String(), "\n"))
}

func TestTaskList(t *testing.T) {
	router := gin.Default()
	c := new(TaskController)
	gin.SetMode(gin.TestMode)
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&models.User{}, &models.Task{})
	db.Create(&models.User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	tasks := []models.Task{
		models.Task{ID: 200, UserID: 100, Title: "task1", Type: 0, CreatedAt: time.Now()},
		models.Task{ID: 201, UserID: 100, Title: "task2", Type: 0, CreatedAt: time.Now()},
	}
	for _, v := range tasks {
		db.Create(&v)
	}
	router.Use(func(c *gin.Context) {
		c.Set("user", 100)
	})
	router.GET("/", c.List)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	created := struct {
		Tasks []models.Task `json:"tasks"`
	}{}
	err := json.Unmarshal([]byte(strings.TrimRight(w.Body.String(), "\n")), &created)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(created.Tasks))
	assert.Equal(t, 200, created.Tasks[0].ID)
	assert.Equal(t, 201, created.Tasks[1].ID)
}

func TestTaskListValidationErrorUserNotFound(t *testing.T) {
	router := gin.Default()
	c := new(TaskController)
	gin.SetMode(gin.TestMode)
	config.Init("test")
	db.Init()
	db := db.GetDB()
	db.AutoMigrate(&models.User{}, &models.Task{})
	router.GET("/", c.List)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 403, w.Code, w.Body.String())
	assert.Equal(t, `{"error":"user not found"}`, strings.TrimRight(w.Body.String(), "\n"))
}

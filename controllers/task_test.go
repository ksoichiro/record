package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ksoichiro/record/db"
	"github.com/ksoichiro/record/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestTaskCreateSuccessfully(t *testing.T) {
	router := gin.Default()
	c := new(TaskController)
	gin.SetMode(gin.TestMode)
	db.InitForTest()
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
	db.InitForTest()
	db := db.GetDB()
	db.AutoMigrate(&models.User{}, &models.Task{})
	db.Create(&models.User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/create", strings.NewReader(`{}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"error":"Key: 'TaskCreateForm.Title' Error:Field validation for 'Title' failed on the 'required' tag"}`, strings.TrimRight(w.Body.String(), "\n"))
}

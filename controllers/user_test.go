package controllers

import (
	"encoding/json"
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

func TestUserCreateSuccessfully(t *testing.T) {
	router := gin.Default()
	c := new(UserController)
	router.POST("/create", c.Create)
	gin.SetMode(gin.TestMode)
	db.InitForTest()
	db := db.GetDB()
	db.AutoMigrate(&models.User{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/create",
		strings.NewReader(`{"name":"foo","password":"test"}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"message":"created"}`, strings.TrimRight(w.Body.String(), "\n"))
}

func TestUserCreateValidationError(t *testing.T) {
	router := gin.Default()
	c := new(UserController)
	router.POST("/create", c.Create)
	gin.SetMode(gin.TestMode)
	db.InitForTest()
	db := db.GetDB()
	db.AutoMigrate(&models.User{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/create",
		strings.NewReader(`{"name":"foo","password":""}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"error":"Key: 'UserCreateForm.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`, strings.TrimRight(w.Body.String(), "\n"))
}

func TestUserLogin(t *testing.T) {
	router := gin.Default()
	c := new(UserController)
	router.POST("/login", c.Login)
	gin.SetMode(gin.TestMode)
	db.InitForTest()
	db := db.GetDB()
	db.AutoMigrate(&models.User{})
	db.Create(&models.User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(`{"name":"foo","password":"test"}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	loginResult := struct {
		Token string `json:"token"`
	}{}
	err := json.Unmarshal([]byte(strings.TrimRight(w.Body.String(), "\n")), &loginResult)
	assert.Nil(t, err)
	assert.NotEmpty(t, loginResult.Token, strings.TrimRight(w.Body.String(), "\n"))
}

package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

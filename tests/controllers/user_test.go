package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/ksoichiro/record/db"
	"github.com/ksoichiro/record/models"
	"github.com/ksoichiro/record/server"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func initDB() {
	testDB, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	db.SetDB(testDB)
}

func TestUserCreateSuccessfully(t *testing.T) {
	router := server.NewRouter()
	gin.SetMode(gin.TestMode)
	initDB()
	db := db.GetDB()
	db.AutoMigrate(&models.User{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user/create",
		strings.NewReader(`{"name":"foo","password":"test"}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"message":"created"}`, w.Body.String())
}

func TestUserCreateValidationError(t *testing.T) {
	router := server.NewRouter()
	gin.SetMode(gin.TestMode)
	initDB()
	db := db.GetDB()
	db.AutoMigrate(&models.User{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user/create",
		strings.NewReader(`{"name":"foo","password":""}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"error":"Key: 'UserCreateForm.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`, w.Body.String())
}

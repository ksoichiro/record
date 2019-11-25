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

func TestRecordCreateSuccessfully(t *testing.T) {
	router := gin.Default()
	c := new(RecordController)
	gin.SetMode(gin.TestMode)
	db.InitForTest()
	db := db.GetDB()
	db.AutoMigrate(&models.User{}, &models.Task{}, &models.Record{})
	db.Create(&models.User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	db.Create(&models.Task{ID: 200, UserID: 100, Title: "task1", Description: "task description", Done: false, Type: 0, Amount: 0, CreatedAt: time.Now()})
	db.Create(&models.Task{ID: 201, UserID: 101, Title: "task2", Description: "task description", Done: false, Type: 0, Amount: 0, CreatedAt: time.Now()})
	db.Create(&models.Task{ID: 202, UserID: 100, Title: "task3", Description: "task description", Done: false, Type: 0, Amount: 0, CreatedAt: time.Now()})
	router.Use(func(c *gin.Context) {
		c.Set("user", 100)
	})
	router.POST("/:date/create", c.Create)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/2019-11-20/create", strings.NewReader(`{"task_id":200}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"message":"created"}`, strings.TrimRight(w.Body.String(), "\n"))
}

func TestRecordCreateValidationError(t *testing.T) {
	router := gin.Default()
	c := new(RecordController)
	router.POST("/:date/create", c.Create)
	gin.SetMode(gin.TestMode)
	db.InitForTest()
	db := db.GetDB()
	db.AutoMigrate(&models.User{}, &models.Task{}, &models.Record{})
	db.Create(&models.User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	db.Create(&models.Task{ID: 200, UserID: 100, Title: "task1", Description: "task description", Done: false, Type: 0, Amount: 0, CreatedAt: time.Now()})
	db.Create(&models.Task{ID: 201, UserID: 101, Title: "task2", Description: "task description", Done: false, Type: 0, Amount: 0, CreatedAt: time.Now()})
	db.Create(&models.Task{ID: 202, UserID: 100, Title: "task3", Description: "task description", Done: false, Type: 0, Amount: 0, CreatedAt: time.Now()})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/2019-11-20/create", strings.NewReader(`{}`))
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"error":"Key: 'RecordCreateForm.TaskID' Error:Field validation for 'TaskID' failed on the 'exists' tag"}`, strings.TrimRight(w.Body.String(), "\n"))
}

func TestRecordListValidationErrorForURI(t *testing.T) {
	router := gin.Default()
	c := new(RecordController)
	gin.SetMode(gin.TestMode)
	db.InitForTest()
	db := db.GetDB()
	db.AutoMigrate(&models.User{}, &models.Task{}, &models.Record{})
	router.Use(func(c *gin.Context) {
		c.Set("user", 100)
	})
	// This error won't happen, since the parameter and routing is defined in the router.
	router.POST("/", c.Create)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code, w.Body.String())
}

func TestRecordList(t *testing.T) {
	router := gin.Default()
	c := new(RecordController)
	gin.SetMode(gin.TestMode)
	db.InitForTest()
	db := db.GetDB()
	db.AutoMigrate(&models.User{}, &models.Task{}, &models.Record{})
	db.Create(&models.User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	db.Create(&models.Task{ID: 200, UserID: 100, Title: "task1", Description: "task description", Done: false, Type: 0, Amount: 0, CreatedAt: time.Now()})
	db.Create(&models.Task{ID: 201, UserID: 101, Title: "task2", Description: "task description", Done: false, Type: 0, Amount: 0, CreatedAt: time.Now()})
	db.Create(&models.Task{ID: 202, UserID: 100, Title: "task3", Description: "task description", Done: false, Type: 0, Amount: 0, CreatedAt: time.Now()})
	db.Create(&models.Record{ID: 300, UserID: 100, TargetDate: mustParse("2006-01-02", "2019-11-19"), TaskID: 200, CreatedAt: time.Now()})
	db.Create(&models.Record{ID: 301, UserID: 100, TargetDate: mustParse("2006-01-02", "2019-11-20"), TaskID: 200, CreatedAt: time.Now()})
	db.Create(&models.Record{ID: 302, UserID: 101, TargetDate: mustParse("2006-01-02", "2019-11-20"), TaskID: 201, CreatedAt: time.Now()})
	db.Create(&models.Record{ID: 303, UserID: 100, TargetDate: mustParse("2006-01-02", "2019-11-21"), TaskID: 200, CreatedAt: time.Now()})
	db.Create(&models.Record{ID: 304, UserID: 100, TargetDate: mustParse("2006-01-02", "2019-11-20"), TaskID: 202, CreatedAt: time.Now()})
	router.Use(func(c *gin.Context) {
		c.Set("user", 100)
	})
	router.GET("/:date", c.List)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/2019-11-20", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code, w.Body.String())
	created := struct {
		Records []models.Record `json:"records"`
	}{}
	err := json.Unmarshal([]byte(strings.TrimRight(w.Body.String(), "\n")), &created)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(created.Records), w.Body.String())
	assert.Equal(t, 301, created.Records[0].ID)
	assert.Equal(t, 304, created.Records[1].ID)
}

func TestRecordListValidationError(t *testing.T) {
	router := gin.Default()
	c := new(RecordController)
	gin.SetMode(gin.TestMode)
	db.InitForTest()
	db := db.GetDB()
	db.AutoMigrate(&models.User{}, &models.Task{}, &models.Record{})
	router.Use(func(c *gin.Context) {
		c.Set("user", 100)
	})
	// This error won't happen, since the parameter and routing is defined in the router.
	router.GET("/", c.List)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code, w.Body.String())
}

func TestRecordListValidationErrorInvalidDateFormat(t *testing.T) {
	router := gin.Default()
	c := new(RecordController)
	gin.SetMode(gin.TestMode)
	db.InitForTest()
	db := db.GetDB()
	db.AutoMigrate(&models.User{}, &models.Task{}, &models.Record{})
	router.Use(func(c *gin.Context) {
		c.Set("user", 100)
	})
	router.GET("/:date", c.List)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/100", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code, w.Body.String())
	assert.Equal(t, `{"error":"invalid date"}`, strings.TrimRight(w.Body.String(), "\n"))
}

func TestRecordListValidationErrorUserNotFound(t *testing.T) {
	router := gin.Default()
	c := new(RecordController)
	gin.SetMode(gin.TestMode)
	db.InitForTest()
	db := db.GetDB()
	db.AutoMigrate(&models.User{}, &models.Task{}, &models.Record{})
	router.GET("/:date", c.List)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/2019-11-20", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 403, w.Code, w.Body.String())
	assert.Equal(t, `{"error":"user not found"}`, strings.TrimRight(w.Body.String(), "\n"))
}

func mustParse(layout string, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

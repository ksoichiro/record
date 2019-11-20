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

func mustParse(layout string, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

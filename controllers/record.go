package controllers

import (
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ksoichiro/record/db"
	"github.com/ksoichiro/record/forms"
	"github.com/ksoichiro/record/models"
)

type RecordController struct{}

func (r RecordController) List(c *gin.Context) {
	d := c.Param("date")
	dateExpr := regexp.MustCompile(`[0-9]{4}-[0-9]{2}-[0-9]{2}`)
	if !dateExpr.MatchString(d) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date"})
		return
	}
	targetDate, _ := time.Parse("2019/04/20", d)
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusOK, gin.H{"error": "user not found"})
		return
	}
	db := db.GetDB()
	records := []models.Record{}
	db.Where("user_id = ? and target_date = ?", userID, targetDate).Find(&records)
	c.JSON(http.StatusOK, gin.H{"records": &records})
}

func (r RecordController) Create(c *gin.Context) {
	d := c.Param("date")
	dateExpr := regexp.MustCompile(`[0-9]{4}-[0-9]{2}-[0-9]{2}`)
	if !dateExpr.MatchString(d) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date"})
		return
	}
	targetDate, err := time.Parse("2006-01-02", d)
	if err != nil {
		panic(err)
	}
	db := db.GetDB()
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusOK, gin.H{"error": "user not found"})
		return
	}
	var json forms.RecordCreateForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tx := db.Begin()
	var user models.User
	tx.Where("id = ?", userID).Find(&user)
	var task models.Task
	var count int
	tx.Where("id = ? and user_id = ?", *json.TaskID, userID).First(&task).Count(&count)
	if count == 0 {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"error": "task not found"})
		return
	}
	tx.Model(&models.Record{}).Where("user_id = ? and target_date = ? and task_id = ?", userID, targetDate, json.TaskID).Count(&count)
	if 0 < count {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "already created"})
		return
	}
	record := models.Record{
		User:       user,
		TargetDate: targetDate,
		Task:       task,
		CreatedAt:  time.Now(),
	}
	if json.Amount == nil {
		record.Amount = 0
	} else {
		record.Amount = *json.Amount
	}
	tx.Create(&record)
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "created"})
}

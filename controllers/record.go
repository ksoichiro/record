package controllers

import (
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ksoichiro/record/forms"
	"github.com/ksoichiro/record/models"
)

// RecordController handles requests about records.
type RecordController struct{}

type recordParam struct {
	Date string `uri:"date" binding:"required"`
}

// List returns the records of the user.
func (r RecordController) List(c *gin.Context) {
	var recordParam recordParam
	if err := c.ShouldBindUri(&recordParam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	d := recordParam.Date
	dateExpr := regexp.MustCompile(`[0-9]{4}-[0-9]{2}-[0-9]{2}`)
	if !dateExpr.MatchString(d) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date"})
		return
	}
	targetDate, _ := time.Parse("2006-01-02", d)
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusOK, gin.H{"error": "user not found"})
		return
	}
	records := models.ListRecords(userID.(int), targetDate)
	c.JSON(http.StatusOK, gin.H{"records": &records})
}

// Create creates a new record of a task for the user.
func (r RecordController) Create(c *gin.Context) {
	var recordParam recordParam
	if err := c.ShouldBindUri(&recordParam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	d := recordParam.Date
	dateExpr := regexp.MustCompile(`[0-9]{4}-[0-9]{2}-[0-9]{2}`)
	if !dateExpr.MatchString(d) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date"})
		return
	}
	var json forms.RecordCreateForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusOK, gin.H{"error": "user not found"})
		return
	}
	if _, err := models.NewRecord(&json, userID.(int), d); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "created"})
}

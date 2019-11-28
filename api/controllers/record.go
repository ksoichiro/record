package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ksoichiro/record/api/forms"
	"github.com/ksoichiro/record/api/models"
)

// RecordController handles requests about records.
type RecordController struct{}

type recordParam struct {
	Date time.Time `uri:"date" time_format:"2006-01-02" time_utc:"1" binding:"required"`
}

// List returns the records of the user.
func (r RecordController) List(c *gin.Context) {
	var recordParam recordParam
	if err := c.ShouldBindUri(&recordParam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "user not found"})
		return
	}
	records := models.ListRecords(userID.(int), recordParam.Date)
	c.JSON(http.StatusOK, gin.H{"records": &records})
}

// Create creates a new record of a task for the user.
func (r RecordController) Create(c *gin.Context) {
	var recordParam recordParam
	if err := c.ShouldBindUri(&recordParam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var json forms.RecordCreateForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "user not found"})
		return
	}
	if _, err := models.NewRecord(&json, userID.(int), recordParam.Date); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "created"})
}

// Update updates the user's record.
func (r RecordController) Update(c *gin.Context) {
	var json forms.RecordUpdateForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "user not found"})
		return
	}
	var record *models.Record
	var err error
	if record, err = models.FindRecord(*json.ID, userID.(int)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := record.Update(&json); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

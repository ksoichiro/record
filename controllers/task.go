package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ksoichiro/record/db"
	"github.com/ksoichiro/record/forms"
	"github.com/ksoichiro/record/models"
)

// TaskController handles requests about tasks.
type TaskController struct{}

// List returns the tasks of the user.
func (t TaskController) List(c *gin.Context) {
	db := db.GetDB()
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "user not found"})
		return
	}
	tasks := []models.Task{}
	db.Where("user_id = ?", userID).Find(&tasks)
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// Create creates a new task for the user.
func (t TaskController) Create(c *gin.Context) {
	var json forms.TaskCreateForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "user not found"})
		return
	}
	models.NewTask(&json, userID.(int))
	c.JSON(http.StatusOK, gin.H{"message": "created"})
}

// Update updates the user's task.
func (t TaskController) Update(c *gin.Context) {
	var json forms.TaskUpdateForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db := db.GetDB()
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "user not found"})
		return
	}
	tx := db.Begin()
	var user models.User
	tx.Where("id = ?", userID).Find(&user)
	var task models.Task
	var count int
	tx.Where("id = ? and user_id = ?", *json.ID, userID).First(&task).Count(&count)
	if count == 0 {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"error": "task not found"})
		return
	}
	if json.Title != nil {
		task.Title = *json.Title
	}
	if json.Description != nil {
		task.Description = *json.Description
	}
	if json.Type != nil {
		task.Type = *json.Type
	}
	if json.Amount != nil {
		task.Amount = *json.Amount
	}
	tx.Save(&task)
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

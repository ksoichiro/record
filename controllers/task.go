package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ksoichiro/record/forms"
	"github.com/ksoichiro/record/models"
)

// TaskController handles requests about tasks.
type TaskController struct{}

// List returns the tasks of the user.
func (t TaskController) List(c *gin.Context) {
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "user not found"})
		return
	}
	tasks := models.ListTasks(userID.(int))
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
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "user not found"})
		return
	}
	var task *models.Task
	var err error
	if task, err = models.FindTask(*json.ID, userID.(int)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := task.Update(&json); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ksoichiro/record/db"
	"github.com/ksoichiro/record/forms"
	"github.com/ksoichiro/record/models"
)

type TaskController struct{}

func (t TaskController) List(c *gin.Context) {
	db := db.GetDB()
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusOK, gin.H{"error": "user not found"})
		return
	}
	tasks := []models.Task{}
	db.Where("user_id = ?", userID).Find(&tasks)
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (t TaskController) Create(c *gin.Context) {
	var json forms.TaskCreateForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db := db.GetDB()
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusOK, gin.H{"error": "user not found"})
		return
	}
	tx := db.Begin()
	var user models.User
	tx.Where("id = ?", userID).Find(&user)
	task := models.Task{
		User:        user,
		Title:       json.Title,
		Description: json.Description,
		CreatedAt:   time.Now(),
	}
	if json.Done == nil {
		task.Done = false
	} else {
		task.Done = *json.Done
	}
	if json.Type == nil {
		task.Type = 0
	} else {
		task.Type = *json.Type
	}
	if json.Amount == nil {
		task.Amount = 0
	} else {
		task.Amount = *json.Amount
	}
	tx.Create(&task)
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "created"})
}

func (t TaskController) Update(c *gin.Context) {
	var json forms.TaskUpdateForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db := db.GetDB()
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusOK, gin.H{"error": "user not found"})
		return
	}
	tx := db.Begin()
	var user models.User
	tx.Where("id = ?", userID).Find(&user)
	var task models.Task
	tx.Where("id = ?", *json.ID).First(&task)
	if json.Title != nil {
		task.Title = *json.Title
	}
	if json.Description != nil {
		task.Description = *json.Description
	}
	if json.Done != nil {
		task.Done = *json.Done
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

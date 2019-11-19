package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ksoichiro/record/db"
	"github.com/ksoichiro/record/models"
)

type TaskController struct{}

func (t TaskController) List(c *gin.Context) {
	db := db.GetDB()
	tasks := []models.Task{}
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusOK, gin.H{"error": "user not found"})
		return
	}
	db.Where("user_id = ?", userID).Find(&tasks)
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (t TaskController) Create(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "created"})
}

package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ksoichiro/record/forms"
	"github.com/ksoichiro/record/models"
)

// UserController handles requests about users.
type UserController struct {
}

// Create creates a new user.
func (u UserController) Create(c *gin.Context) {
	var json forms.UserCreateForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := models.NewUser(&json)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "created"})
}

// Login authenticates the user and returns a authentication token.
func (u UserController) Login(c *gin.Context) {
	var json forms.UserLoginForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := models.Login(&json)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

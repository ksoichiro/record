package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TaskController struct{}

func (t TaskController) TaskCreate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "created"})
}

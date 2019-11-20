package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ksoichiro/record/controllers"
	"github.com/ksoichiro/record/db"
	"github.com/ksoichiro/record/middlewares"
	"github.com/ksoichiro/record/models"
)

func adminUser(c *gin.Context) {
	db := db.GetDB()
	users := []models.User{}
	db.Find(&users)
	c.JSON(http.StatusOK, gin.H{"users": users})
}

// NewRouter creates and sets up a new router for 'gin'.
func NewRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		db := db.GetDB()
		task := models.Task{}
		db.First(&task)
		fmt.Printf("%d: UserId: %d, Title: %s\n", task.ID, task.UserID, task.Title)

		c.String(http.StatusOK, "pong")
	})

	userGroup := r.Group("/user")
	user := new(controllers.UserController)
	userGroup.POST("/create", user.Create)
	userGroup.POST("/login", user.Login)

	taskGroup := r.Group("/task")
	taskGroup.Use(middlewares.UserAuthenticator())
	task := new(controllers.TaskController)
	taskGroup.GET("", task.List)
	taskGroup.POST("/create", task.Create)
	taskGroup.POST("/update", task.Update)

	recordGroup := r.Group("/record")
	recordGroup.Use(middlewares.UserAuthenticator())
	record := new(controllers.RecordController)
	recordGroup.GET("/:date", record.List)
	recordGroup.POST("/:date/create", record.Create)

	adminGroup := r.Group("/admin")
	adminGroup.GET("/user", adminUser)

	return r
}

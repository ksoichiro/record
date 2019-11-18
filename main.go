package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type userCreateForm struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type task struct {
	ID     int
	UserID int
	Title  string
}

func connect() *gorm.DB {
	db, err := gorm.Open("mysql", "app:app@tcp(127.0.0.1:3306)/test?parseTime=true&loc=Asia%2FTokyo")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		db := connect()
		defer db.Close()
		task := task{}
		db.First(&task)
		fmt.Printf("%d: UserId: %d, Title: %s\n", task.ID, task.UserID, task.Title)

		c.String(http.StatusOK, "pong")
	})

	r.POST("/user/create", func(c *gin.Context) {
		var json userCreateForm
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		name := json.Name
		hash, err := bcrypt.GenerateFromPassword([]byte(json.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		password := string(hash)
		fmt.Printf("%s / %s\n", name, password)
		db := connect()
		defer db.Close()
		tx := db.Begin()
		user := user{Name: name, Password: password, CreatedAt: time.Now()}
		tx.Create(&user)
		tx.Commit()

		c.JSON(http.StatusOK, gin.H{"message": "created"})
	})

	r.GET("/admin/user", func(c *gin.Context) {
		db := connect()
		defer db.Close()
		users := []user{}
		db.Find(&users)
		c.JSON(http.StatusOK, gin.H{"users": users})
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

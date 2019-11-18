package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type user struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
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
		name := c.PostForm("name")
		password := c.PostForm("password")
		fmt.Printf("%s / %s\n", name, password)
		db := connect()
		defer db.Close()
		tx := db.Begin()
		user := user{Name: name, CreatedAt: time.Now()}
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

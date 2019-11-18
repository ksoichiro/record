package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type task struct {
	ID    int
	Title string
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		db, err := sql.Open("mysql", "app:app@tcp(127.0.0.1:3306)/test")
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		defer db.Close()
		rows, err := db.Query("SELECT id, title FROM task")
		if err != nil {
			fmt.Printf("%v\n", err)
			c.String(http.StatusInternalServerError, "failed")
			return
		}
		for rows.Next() {
			task := task{}
			err = rows.Scan(&task.ID, &task.Title)
			fmt.Printf("%d: %s\n", task.ID, task.Title)
		}

		c.String(http.StatusOK, "pong")
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

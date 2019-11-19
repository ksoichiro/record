package models

type Task struct {
	ID     int    `json:"id"`
	UserID int    `json:"-"`
	User   User   `json:"-"`
	Title  string `json:"title"`
}

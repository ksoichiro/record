package main

import (
	"github.com/ksoichiro/record/db"
	"github.com/ksoichiro/record/server"
)

func main() {
	db.Init()
	defer db.CloseDB()
	server.Init()
}

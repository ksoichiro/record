package main

import (
	"flag"

	"github.com/ksoichiro/record/api/config"
	"github.com/ksoichiro/record/api/db"
	"github.com/ksoichiro/record/api/server"
)

func main() {
	environment := flag.String("e", "development", "")
	flag.Parse()
	config.Init(*environment)
	db.Init()
	defer db.CloseDB()
	server.Init()
}

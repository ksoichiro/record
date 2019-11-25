package main

import (
	"flag"

	"github.com/ksoichiro/record/config"
	"github.com/ksoichiro/record/db"
	"github.com/ksoichiro/record/server"
)

func main() {
	environment := flag.String("e", "development", "")
	flag.Parse()
	config.Init(*environment)
	db.Init()
	defer db.CloseDB()
	server.Init()
}

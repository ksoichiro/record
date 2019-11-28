package server

import "github.com/ksoichiro/record/api/config"

// Init initializes the server and runs it.
func Init() {
	r := NewRouter()
	r.Run(":" + config.GetConfig().GetString("server.port"))
}

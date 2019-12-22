package server

import "github.com/ksoichiro/record/api/config"

// Init initializes the server and runs it.
func Init() {
	r := newRouter()
	_ = r.Run(":" + config.GetConfig().GetString("server.port"))
}

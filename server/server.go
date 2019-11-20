package server

// Init initializes the server and runs it.
func Init() {
	r := NewRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

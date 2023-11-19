package main

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/michaelmagen/task-together/db"
	"github.com/michaelmagen/task-together/internal/configs"
)

// Things to make:
// 1. Auth (google oauth + redis)
// 2. db for storing data
// 3. websockets for live updates
// 4. routes for users, lists, and tasks
func main() {
	// TODO: Use viper for env variables
	// Initialize Viper across the application
	configs.InitializeViper()

	// Initialize database
	db.InitDB()

	// TODO: Add melody websocket server
	r := chi.NewRouter()

	// Middleware
	// TODO: Create auth middleware
	r.Use(middleware.Logger)

	// Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from task together server"))
	})

	// TODO: Create auth route

	http.ListenAndServe(":3000", r)
}

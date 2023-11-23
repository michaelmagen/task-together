package main

import (
	"net/http"

	"encoding/gob"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/michaelmagen/task-together/auth"
	"github.com/michaelmagen/task-together/configs"
	"github.com/michaelmagen/task-together/model"
	"golang.org/x/oauth2"
)

// TODO: figure out diff between r.context and context.Background()

// Things to make:
// 1. Auth (google oauth + redis)
// 2. db for storing data
// 3. websockets for live updates
// 4. routes for users, lists, and tasks

func init() {
	// Initialize Viper across the application
	configs.InitializeViper()

	// Initialize database
	model.InitDB()

	// Initialize session store
	configs.InitSessionStore()

	// Initialize google oauth client
	auth.InitializeOAuthGoogle()

	// Add encoding for oauth tokens to session store
	gob.Register(&oauth2.Token{})
}

func main() {
	// TODO: Add melody websocket server
	r := chi.NewRouter()

	// Middleware
	// TODO: Create auth middleware
	r.Use(middleware.Logger)

	// Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	r.Get("/done", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("in protected route"))
	})

	r.Get("/login", auth.HandleGoogleLogin)
	r.Get("/auth/callback", auth.CallbackGoogleOauth)

	// TODO: Create auth route

	http.ListenAndServe(":3000", r)
}

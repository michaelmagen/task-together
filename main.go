package main

import (
	"encoding/gob"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/michaelmagen/task-together/auth"
	"github.com/michaelmagen/task-together/configs"
	"github.com/michaelmagen/task-together/model"
	"github.com/michaelmagen/task-together/routes"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

// Things to make:
// 3. websockets for live updates
// 4. routes for users, lists, and tasks, and invitations

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
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{viper.GetString("frontendURL")},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "PATCH", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"X-Requested-With", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Routes
	r.With(auth.AuthedWithGoogle).Route("/users", routes.UsersRoute)
	r.With(auth.AuthedWithGoogle).Route("/lists", routes.ListsRoute)
	r.With(auth.AuthedWithGoogle).Route("/invitations", routes.InvitationsRoute)
	r.Get("/login", auth.HandleGoogleLogin)
	r.Get("/auth/callback", auth.CallbackGoogleOauth)

	http.ListenAndServe(viper.GetString("PORT"), r)
}

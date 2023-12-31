package routes

import (
	"encoding/json"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/michaelmagen/task-together/model"
)

func UsersRoute(r chi.Router) {
	r.Get("/", getLoggedInUser)
	r.Get("/{userID}", getUserHandler)
}

func getLoggedInUser(w http.ResponseWriter, r *http.Request) {
	//Extract userID from the context
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		log.Error("getLoggedInUser:", "err", err)
		http.Error(w, "Could not get user", http.StatusInternalServerError)
		return
	}

	// Get user object from db
	user, err := model.GetUserByID(userID)
	if err != nil {
		log.Error("getLoggedInUser:", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set the Content-Type header to indicate JSON
	w.Header().Set("Content-Type", "application/json")

	// Use json.NewEncoder to directly encode and write to the response writer
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Error("getLoggedInUser: Failed to encode user", "err", err)
		http.Error(w, "Could not get lists", http.StatusInternalServerError)
		return
	}

}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get userID from param
	userIDParam := chi.URLParam(r, "userID")
	if userIDParam == "" {
		log.Error("getUserHandler: no userID in url param")
		http.Error(w, "userID missing in url", http.StatusBadRequest)
		return
	}
	userID := model.UserID(userIDParam)

	// Get user object from db
	user, err := model.GetUserByID(userID)
	if err != nil {
		log.Error("getUserHandler", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set the Content-Type header to indicate JSON
	w.Header().Set("Content-Type", "application/json")

	// Use json.NewEncoder to directly encode and write to the response writer
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Error("getUserHandler: Failed to encode user", "err", err)
		http.Error(w, "Could not get lists", http.StatusInternalServerError)
		return
	}
}

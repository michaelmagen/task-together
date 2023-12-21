package routes

import (
	"encoding/json"
	"net/http"

	"fmt"
	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/michaelmagen/task-together/model"
)

func ListsRoute(r chi.Router) {
	r.Get("/", userListsHandler)
	r.Post("/", postListHandler)
	r.Route("/{listID}/tasks", TasksRoute)
}

func userListsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract userID from the context
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		log.Error("/list/all:", "err", err)
		http.Error(w, "Could not get lists", http.StatusInternalServerError)
		return
	}

	lists, err := model.GetAllListForUser(userID)
	if err != nil {
		log.Error("Failed to get all lists for user", "err", err)
		http.Error(w, "Could not get lists", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to indicate JSON
	w.Header().Set("Content-Type", "application/json")

	// Use json.NewEncoder to directly encode and write to the response writer
	if err := json.NewEncoder(w).Encode(lists); err != nil {
		log.Error("Failed to encode lists", "err", err)
		http.Error(w, "Could not get lists", http.StatusInternalServerError)
		return
	}
}

type postListRequest struct {
	Name string `json:"name"`
}

func postListHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		log.Error("postListHandler:", "err", err)
		http.Error(w, "Could not get lists", http.StatusInternalServerError)
		return
	}

	// Get list name from request
	var requestBody postListRequest
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Error("postListHandler:", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Add list to db
	listID, err := model.CreateList(requestBody.Name, userID)
	if err != nil {
		log.Error("postListHandler:", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Respond with the new listID
	response := map[string]string{"listID": fmt.Sprintf("%d", listID)} // Convert listID to string
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/michaelmagen/task-together/model"
)

func TasksRoute(r chi.Router) {
	r.Get("/", getTasksInListHandler)
	r.Post("/", createTaskHandler)
	r.Delete("/{taskID}", deleteTaskHandler)
	r.Patch("/{taskID}", toggleTaskCompletionHandler)
}

func getTasksInListHandler(w http.ResponseWriter, r *http.Request) {
	// Get list id param and convert it to int
	listIDParam, err := strconv.Atoi(chi.URLParam(r, "listID"))
	if err != nil {
		log.Error("getTasksInListHandler:", "err", err)
		http.Error(w, "Invalid listID", http.StatusBadRequest)
		return
	}
	listID := model.ListID(listIDParam)

	// Retrieve tasks from db
	tasks, err := model.GetTasksByList(listID)
	if err != nil {
		log.Error("getTasksInListHandler:", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Encode tasks as json and send as response
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		log.Error("getTasksInList: Failed to encode tasks", "err", err)
		http.Error(w, "Failed to encode tasks", http.StatusInternalServerError)
		return
	}
}

type createTaskRequst struct {
	Content string `json:"content"`
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Get user id
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		log.Error("createTaskHandler:", "err", err)
		http.Error(w, "Failed to get user id from request", http.StatusBadRequest)
		return
	}

	// Get list id param and convert it to int
	listIDParam, err := strconv.Atoi(chi.URLParam(r, "listID"))
	if err != nil {
		log.Error("createTaskHandler:", "err", err)
		http.Error(w, "Invalid listID", http.StatusBadRequest)
		return
	}
	listID := model.ListID(listIDParam)

	// Get the content for the task
	var requestBody createTaskRequst
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Error("createTaskHandler:", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	content := requestBody.Content

	// Create task in db
	task, err := model.CreateTask(content, listID, userID)
	if err != nil {
		log.Error("createTaskHandler:", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return created task
	w.Header().Set("Content-Type", "application/json")

	// TODO: Might be better just to return the task id here instead of the whole task
	// Encode task as json and send as response
	if err := json.NewEncoder(w).Encode(task); err != nil {
		log.Error("Failed to encode task", "err", err)
		http.Error(w, "Could not create task", http.StatusInternalServerError)
		return
	}
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Get task id param and convert it to int
	taskIDParam, err := strconv.Atoi(chi.URLParam(r, "taskID"))
	if err != nil {
		log.Error("deleteTaskHandler:", "err", err)
		http.Error(w, "Invalid taskID", http.StatusBadRequest)
		return
	}
	taskID := model.TaskID(taskIDParam)

	err = model.DeleteTask(taskID)
	if err != nil {
		log.Error("deleteTaskHandler:", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Send the deletion event through websockets

	// Successful deletion, send a 204 No Content response
	w.WriteHeader(http.StatusNoContent)
}

type toggleTaskRequestBody struct {
	Completed string `json:"completed"`
}

func toggleTaskCompletionHandler(w http.ResponseWriter, r *http.Request) {
	// Get user id
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		log.Error("toggleTaskCompletionHandler:", "err", err)
		http.Error(w, "Failed to get user id from request", http.StatusBadRequest)
		return
	}

	// Get task id param and convert it to int
	taskIDParam, err := strconv.Atoi(chi.URLParam(r, "taskID"))
	if err != nil {
		log.Error("deleteTaskHandler:", "err", err)
		http.Error(w, "Invalid taskID", http.StatusBadRequest)
		return
	}
	taskID := model.TaskID(taskIDParam)

	// Get the completion status for the task
	var requestBody toggleTaskRequestBody
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Error("toogleTaskCompletionHandler:", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert completion into boolean
	completed, err := strconv.ParseBool(requestBody.Completed)
	if err != nil {
		log.Error("toogleTaskCompletionHandler:", "err", err)
		http.Error(w, "completed must be a boolean", http.StatusBadRequest)
		return
	}

	err = model.UpdateTaskCompletion(completed, taskID, userID)
	if err != nil {
		log.Error("togleTaskCompletionHandler:", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Send the new updated task through the websocket
	// Successful deletion, send a 204 No Content response
	w.WriteHeader(http.StatusNoContent)
}

package routes

import (
	"errors"
	"net/http"

	"github.com/michaelmagen/task-together/model"
)

func getUserIDFromRequest(r *http.Request) (model.UserID, error) {
	// Extract userID from the context
	userIDFromContext, ok := r.Context().Value("userID").(string)

	if !ok || userIDFromContext == "" {
		return model.UserID(""), errors.New("context: Failed to get userID from context")
	}

	userID := model.UserID(userIDFromContext)
	return userID, nil
}

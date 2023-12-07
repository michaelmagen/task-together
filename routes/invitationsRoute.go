package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/michaelmagen/task-together/model"
)

func InvitationsRoute(r chi.Router) {
	r.Get("/", getUserInvitationsHandler)
	r.Post("/", createInvitationHandler)
	r.Delete("/{invitationID}", deleteInvitationHandler)
}

func getUserInvitationsHandler(w http.ResponseWriter, r *http.Request) {
	// Get user id
	userID, err := getUserIDFromRequest(r)
	if err != nil {
		log.Error("getUserInvitationHandler:", "err", err)
		http.Error(w, "Failed to get user id from request", http.StatusBadRequest)
		return
	}

	// Get invitations from DB
	invitations, err := model.GetInvitationsByReceiver(userID)
	if err != nil {
		log.Error("getUserInvitationHandler:", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set the Content-Type header to indicate JSON
	w.Header().Set("Content-Type", "application/json")

	// Use json.NewEncoder to directly encode and write to the response writer
	if err := json.NewEncoder(w).Encode(invitations); err != nil {
		log.Error("getUserInvitationsHandler: Failed to encode", "err", err)
		http.Error(w, "Could not get lists", http.StatusInternalServerError)
		return
	}
}

type createInviteRequestBody struct {
	ReceiverEmail string       `json:"receiver_email"`
	ListID        model.ListID `json:"list_id"`
}

func createInvitationHandler(w http.ResponseWriter, r *http.Request) {
	// Get id of user that make request, they are the sender
	senderID, err := getUserIDFromRequest(r)
	if err != nil {
		log.Error("getUserInvitationHandler:", "err", err)
		http.Error(w, "Failed to get user id from request", http.StatusBadRequest)
		return
	}
	// Get receiver email and list id from request body
	var requestBody createInviteRequestBody
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Error("createInvitationHandler:", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	receiverEmail := requestBody.ReceiverEmail
	listID := requestBody.ListID

	err = model.CreateInvitation(senderID, receiverEmail, listID)
	if err != nil {
		log.Error("createInvitatinoHandler:", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: use websocket to send invite in real time

	// Successful deletion, send a 204 No Content response
	w.WriteHeader(http.StatusNoContent)
}

// We delete the invite. If accepted is true, then we delete the invite and add the user to the list they were invited to.
type deleteInvitationRequestBody struct {
	Accepted string       `json:"accepted"`
	ListID   model.ListID `json:"list_id"`
}

func deleteInvitationHandler(w http.ResponseWriter, r *http.Request) {
	// Get invitation id param and convert it to int
	invitationIDParam, err := strconv.Atoi(chi.URLParam(r, "invitationID"))
	if err != nil {
		log.Error("invitationIDParam:", "err", err)
		http.Error(w, "Invalid invitationID", http.StatusBadRequest)
		return
	}
	invitationID := model.InvitationID(invitationIDParam)

	// Get the acceptance status for the invite
	var requestBody deleteInvitationRequestBody
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Error("deleteInvitationHandler:", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	listID := requestBody.ListID

	// Convert accepted into boolean
	acceptedInvite, err := strconv.ParseBool(requestBody.Accepted)
	if err != nil {
		log.Error("deleteInvitationHandler:", "err", err)
		http.Error(w, "accepted must be a boolean", http.StatusBadRequest)
		return
	}

	// If the invite was accepted, then we want to add the user to the the list
	if acceptedInvite {
		// Get user id
		userID, err := getUserIDFromRequest(r)
		if err != nil {
			log.Error("deleteInvitationHandler:", "err", err)
			http.Error(w, "Failed to get user id from request", http.StatusBadRequest)
			return
		}
		err = model.AddUserToList(listID, userID)
		if err != nil {
			log.Error("deleteInvitationHandler:", "err", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Delete the invitation
	err = model.DeleteInvitation(invitationID)
	if err != nil {
		log.Error("deleteInvitationHandler:", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// If invite accepted, send back the list. Otherwise send no content
	if acceptedInvite {
		list, err := model.GetListByID(listID)
		if err != nil {
			log.Error("deleteInvitationHandler:", "err", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Set the Content-Type header to indicate JSON
		w.Header().Set("Content-Type", "application/json")

		// Use json.NewEncoder to directly encode and write to the response writer
		if err := json.NewEncoder(w).Encode(list); err != nil {
			log.Error("deletListHandler: Failed to encode", "err", err)
			http.Error(w, "Could return list", http.StatusInternalServerError)
			return
		}
	} else {
		// Successful deletion, send a 204 No Content response
		w.WriteHeader(http.StatusNoContent)
	}

}

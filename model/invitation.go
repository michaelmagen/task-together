package model

import (
	"github.com/charmbracelet/log"
)

type InvitationID int

type Invitation struct {
	InvitationID InvitationID `json:"invitation_id"`
	SenderID     UserID       `json:"sender_id"`
	ReceiverID   UserID       `json:"receiver_id"`
	ListID       ListID       `json:"list_id"`
}

func CreateInvitation(senderID UserID, receiverID UserID, listID ListID) (InvitationID, error) {
	statement := `
		INSERT INTO invitations (sender_id, receiver_id, list_id)
		VALUES ($1, $2, $3)
		RETURNING invitation_id
	`

	var invitationID InvitationID
	err := db.QueryRow(statement, senderID, receiverID, listID).Scan(&invitationID)
	if err != nil {
		log.Error("Fail to create invitation", "err", err)
		return 0, nil
	}

	return invitationID, nil
}

func DeleteInvitation(invitationID InvitationID) error {
	statement := "DELETE FROM invitations WHERE invitation_id = $1"

	_, err := db.Exec(statement, invitationID)
	if err != nil {
		log.Error("Fail to delete invitation", "err", err)
		return err
	}
	return nil
}

// Deletes the invitation from the db and add the user to the list
func AcceptInvitation(invitationID InvitationID) error {
	// First delete the invitation, while getting the userID that need to add to listID
	statement := `
		DELETE FROM invitations
		WHERE invitation_id = $1
		RETURNING receiver_id, list_id
	`

	var receiverID UserID
	var listID ListID
	err := db.QueryRow(statement, invitationID).Scan(&receiverID, &listID)
	if err != nil {
		log.Error("Fail to delete while accepting", "err", err)
		return err
	}
	// Add the user to the list
	err = AddUserToList(listID, receiverID)
	if err != nil {
		log.Error("Fail to add user to list while accepting invite", "err", err)
		return err
	}

	return nil
}

func GetInvitationsByUserID(userID UserID) ([]Invitation, error) {
	statement := `
		SELECT invitation_id, sender_id, receiver_id, list_id
		FROM invitations
		WHERE receiver_id = $1
	`

	rows, err := db.Query(statement, userID)
	if err != nil {
		log.Error("Fail to get invitations by user ID", "err", err)
		return nil, err
	}
	defer rows.Close()

	var invitations []Invitation

	for rows.Next() {
		var invitation Invitation
		err := rows.Scan(&invitation.InvitationID, &invitation.SenderID, &invitation.ReceiverID, &invitation.ListID)
		if err != nil {
			log.Error("Error scanning row (GetInvitationByUserID):", "err", err)
			return nil, err
		}
		invitations = append(invitations, invitation)
	}

	if err := rows.Err(); err != nil {
		log.Error("Error iterating over rows:", "err", err)
		return nil, err
	}

	return invitations, nil
}

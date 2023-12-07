package model

import (
	"github.com/charmbracelet/log"
)

type InvitationID int

type Invitation struct {
	InvitationID InvitationID `json:"invitation_id"`
	ReceiverID   UserID       `json:"receiver_id"`
	Sender       *User        `json:"sender"`
	List         *List        `json:"list"`
}

// Do not need to return created invite since this only needs to be gotten by the reciever not sender of invite
func CreateInvitation(senderID UserID, receiverEmail string, listID ListID) error {
	// Get userID for receiver
	receiverID, err := getUserIDFromEmail(receiverEmail)
	if err != nil {
		log.Error("CreateInvitation: Failed to get receiverID", "err", err)
		return err
	}

	statement := `
		INSERT INTO invitations (sender_id, receiver_id, list_id)
		VALUES ($1, $2, $3)
	`

	_, err = db.Exec(statement, senderID, receiverID, listID)
	if err != nil {
		log.Error("Fail to create invitation", "err", err)
		return err
	}

	return nil
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

func GetInvitationsByReceiver(userID UserID) ([]Invitation, error) {
	// Get all invitations sent to userID, join it with info about sender and list
	statement := `
		SELECT 
			i.invitation_id, 
			i.receiver_id, 
			u.sender_id,
		    u.email, 
			u.verified_email, 
			u.name, 
			u.given_name, 
		    u.family_name, 
			u.picture, 
			u.locale,
			l.list_id, 
			l.name, 
			l.creator_id, 
			l.created_at,
		FROM invitations i
		INNER JOIN users u ON i.sender_id = u.user_id
		INNER JOIN lists l ON i.list_id = l.list_id
		WHERE i.receiver_id = $1
	`

	// Run query
	rows, err := db.Query(statement, userID)
	if err != nil {
		log.Error("Fail to get invitations by user ID", "err", err)
		return nil, err
	}
	defer rows.Close()

	var invitations []Invitation

	for rows.Next() {
		var invitation Invitation
		var sender User
		var list List
		err := rows.Scan(
			&invitation.InvitationID,
			&invitation.ReceiverID,
			&sender.UserID,
			&sender.Email,
			&sender.VerifiedEmail,
			&sender.Name,
			&sender.GivenName,
			&sender.FamilyName,
			&sender.Picture,
			&sender.Locale,
			&list.ListID,
			&list.Name,
			&list.CreatorID,
			&list.CreatedAt,
		)
		if err != nil {
			log.Error("Error scanning row (GetInvitationByUserID):", "err", err)
			return nil, err
		}
		// Add Sender and List object to invitation
		invitation.Sender = &sender
		invitation.List = &list
		invitations = append(invitations, invitation)
	}

	// Check for error while iterating through rows
	if err := rows.Err(); err != nil {
		log.Error("Error iterating over rows(GetInvitationByUserID):", "err", err)
		return nil, err
	}

	return invitations, nil
}

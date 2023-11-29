package model

import (
	"github.com/charmbracelet/log"
	"time"
)

type ListID int

// List represents the lists table in the database
type List struct {
	ListID    ListID    `json:"list_id"`
	Name      string    `json:"name"`
	CreatorID UserID    `json:"creator_id"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateList(listName string, userID UserID) (ListID, error) {
	var listID ListID

	// Create list in database
	err := db.QueryRow(`
		INSERT INTO lists (name, creator_id)
		VALUES ($1, $2)
		RETURNING list_id
	`, listName, userID).Scan(&listID)

	if err != nil {
		log.Error("Failed to create list in db", "err", err)
		return 0, nil
	}

	// Make sure user is added to list through user_list table
	err = AddUserToList(listID, userID)
	if err != nil {
		return 0, nil
	}
	return listID, nil
}

func GetAllListForUser(userID UserID) ([]List, error) {
	rows, err := db.Query(`
		SELECT l.list_id, l.name, l.creator_id, l.created_at
		FROM lists l
		INNER JOIN user_list ul ON l.list_id = ul.list_id
		WHERE ul.user_id = $1
	`, userID)

	if err != nil {
		log.Error("Error querying the database:", "err", err)
		return nil, err
	}
	defer rows.Close()

	var lists []List

	for rows.Next() {
		var list List
		err := rows.Scan(&list.ListID, &list.Name, &list.CreatorID, &list.CreatedAt)
		if err != nil {
			log.Error("Error scanning row:", "err", err)
			return nil, err
		}
		lists = append(lists, list)
	}

	if err := rows.Err(); err != nil {
		log.Error("Error iterating over rows:", "err", err)
		return nil, err
	}

	return lists, nil
}

func AddUserToList(listID ListID, userID UserID) error {
	// Add list and user combo to user_list table
	_, err := db.Exec(`
		INSERT INTO user_list (user_id, list_id)
		VALUES ($1, $2)
	`, userID, listID)

	if err != nil {
		log.Error("Failed to insert user list combo", "err", err)
		return err
	}

	return nil
}

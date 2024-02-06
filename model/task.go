package model

import (
	"github.com/charmbracelet/log"
	"time"
)

type TaskID int

type Task struct {
	TaskID      TaskID         `json:"task_id"`
	Content     string         `json:"content"`
	ListID      ListID         `json:"list_id"`
	CreatorID   UserID         `json:"creator_id"`
	CompleterID NullableUserID `json:"completer_id,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	Completed   bool           `json:"completed"`
	Creator     *User          `json:"creator"`
	Completer   *User          `json:"completer,omitempty"`
}

func CreateTask(content string, listID ListID, creatorID UserID) (*Task, error) {
	statement := `
		INSERT INTO	tasks (content, list_id, creator_id)
		VALUES ($1, $2, $3)
		RRETURNING task_id, content, list_id, creator_id, created_at, completed, completer_id
	`
	var task Task
	err := db.QueryRow(statement, content, listID, creatorID).Scan(
		&task.TaskID, &task.Content, &task.ListID, &task.CreatorID, &task.CreatedAt, &task.Completed, &task.CompleterID,
	)
	if err != nil {
		log.Error("Fail to create task", "err", err)
		return nil, err
	}

	err = attachCreatorAndCompleterInfo(&task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func DeleteTask(taskID TaskID) error {

	statement := "DELETE FROM tasks WHERE task_id = $1"

	_, err := db.Exec(statement, taskID)
	if err != nil {
		log.Error("Fail to delete task", "err", err)
		return err
	}
	return nil
}

func GetTasksByList(listID ListID) ([]Task, error) {

	statement := `
		SELECT task_id, content, list_id, creator_id, created_at, completed, completer_id
		FROM tasks
		WHERE list_id = $1
	`

	rows, err := db.Query(statement, listID)
	if err != nil {
		log.Error("Fail to get tasks", "err", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.TaskID, &task.Content, &task.ListID, &task.CreatorID, &task.CreatedAt, &task.Completed, &task.CompleterID)
		if err != nil {
			log.Error("Failed to scan task", "err", err)
			return nil, err
		}
		err = attachCreatorAndCompleterInfo(&task)
		if err != nil {
			log.Error("GetTasksByList:", "err", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		log.Error("Error iterating over rows (GetTasksByList):", "err", err)
		return nil, err
	}

	return tasks, nil
}

func UpdateTaskCompletion(completed bool, taskID TaskID, userID UserID) error {
	if completed == true {
		return completeTask(taskID, userID)
	} else {
		return uncompleteTask(taskID)
	}
}

func completeTask(taskID TaskID, userID UserID) error {
	statement := `
		UPDATE tasks
		SET completed = true, completer_id = $1
		WHERE task_id = $2
	`

	_, err := db.Exec(statement, userID, taskID)
	if err != nil {
		log.Error("Failed to complete task", "err", err)
		return err
	}

	return nil
}

func uncompleteTask(taskID TaskID) error {
	statement := `
		UPDATE tasks
		SET completed = false, completer_id = NULL 
		WHERE task_id = $1 
	`

	_, err := db.Exec(statement, taskID)
	if err != nil {
		log.Error("Failed to uncomplete task", "err", err)
		return err
	}

	return nil
}

// TODO: Make this one sql operation through joins
func attachCreatorAndCompleterInfo(task *Task) error {
	var err error
	// Fetch creator information
	task.Creator, err = GetUserByID(task.CreatorID)
	if err != nil {
		log.Error("Fail to fetch creator information", "err", err)
		return err
	}

	// If task has a completer, fetch completer information
	if task.CompleterID.Valid {
		task.Completer, err = GetUserByID(task.CompleterID.GetID())
		if err != nil {
			log.Error("Fail to fetch completer information", "err", err)
			return err
		}
	}
	return nil
}

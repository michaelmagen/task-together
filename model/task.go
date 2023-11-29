package model

import (
	"github.com/charmbracelet/log"
	"time"
)

type TaskID int

type Task struct {
	TaskID    TaskID    `json:"task_id"`
	Content   string    `json:"content"`
	ListID    ListID    `json:"list_id"`
	CreatorID UserID    `json:"creator_id"`
	CreatedAt time.Time `json:"created_at"`
	Completed bool      `json:"completed"`
}

func CreateTask(content string, listID ListID, creatorID UserID) (TaskID, error) {
	statement := `
		INSERT INTO	tasks (content, list_id, creator_id)
		VALUES ($1, $2, $3)
		RETURNING task_id
	`
	var taskID TaskID
	err := db.QueryRow(statement, content, listID, creatorID).Scan(&taskID)
	if err != nil {
		log.Error("Fail to create task", "err", err)
		return 0, err
	}

	return taskID, nil
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

func GetTasksByList(listID int) ([]Task, error) {

	statement := `
		SELECT task_id, content, list_id, creator_id, created_at, completed
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
		err := rows.Scan(&task.TaskID, &task.Content, &task.ListID, &task.CreatorID, &task.CreatedAt, &task.Completed)
		if err != nil {
			log.Error("Failed to scan task", "err", err)
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

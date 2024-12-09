package database

import (
	"database/sql"
	"github.com/GoSPB/go_final/internal/models"
)

// GetTask возвращает задачу по её ID
func GetTask(db *sql.DB, id string) (models.Task, error) {
	var task models.Task
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return task, err
	}
	return task, nil
}

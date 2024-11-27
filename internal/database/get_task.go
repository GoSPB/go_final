package database

import (
	"database/sql"

	"github.com/anton-ag/todolist/internal/models"
)

func GetTask(db *sql.DB, id string) (models.Task, error) {
	var task models.Task

	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id"
	row := db.QueryRow(
		query,
		sql.Named("id", id),
	)
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return models.Task{}, err
	}
	return task, nil
}

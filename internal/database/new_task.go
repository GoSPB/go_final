package database

import (
	"database/sql"

	"github.com/anton-ag/todolist/internal/models"
)

func NewTask(db *sql.DB, task models.Task) (int64, error) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)"
	res, err := db.Exec(
		query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

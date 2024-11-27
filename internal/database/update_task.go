package database

import (
	"database/sql"

	"github.com/anton-ag/todolist/internal/models"
)

func UpdateTask(db *sql.DB, task models.Task) (int64, error) {
	query := "UPDATE scheduler SET id = :id, date = :date, title = :title, comment= :comment, repeat= :repeat WHERE id = :id"
	res, err := db.Exec(
		query,
		sql.Named("id", task.ID),
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

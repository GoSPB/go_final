package database

import (
	"database/sql"
	"github.com/GoSPB/go_final/internal/models"
)

// UpdateTask обновляет параметры задачи в базе данных
func UpdateTask(db *sql.DB, task models.Task) error {
	query := `
		UPDATE scheduler 
		SET date = ?, title = ?, comment = ?, repeat = ? 
		WHERE id = ?`
	_, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	return err
}

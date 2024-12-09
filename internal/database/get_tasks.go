package database

import (
	"database/sql"
	"github.com/GoSPB/go_final/internal/models"
)

func GetTasks(db *sql.DB) ([]models.Task, error) {
	tasks := []models.Task{}
	query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date"
	rows, err := db.Query(query)
	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

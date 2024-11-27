package database

import (
	"database/sql"
	"strings"
	"time"

	"github.com/anton-ag/todolist/internal/models"
)

func GetTasks(db *sql.DB, search string) ([]models.Task, error) {
	var task models.Task
	var tasks []models.Task
	var searchByDate string

	searchByDate = ""
	dateParsed, err := time.Parse(models.SearchDateFormat, search)
	if err == nil {
		searchByDate = dateParsed.Format(models.DateFormat)
	}

	search = strings.Join([]string{"%", search, "%"}, "")

	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE :search OR comment LIKE :search OR date LIKE :datesearch ORDER BY date LIMIT :limit"
	rows, err := db.Query(
		query,
		sql.Named("limit", models.Limit),
		sql.Named("search", search),
		sql.Named("datesearch", searchByDate),
	)
	defer rows.Close()

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if len(tasks) == 0 {
		return []models.Task{}, err
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

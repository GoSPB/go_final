package repository

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/anton-ag/todolist/internal/models"
	_ "modernc.org/sqlite"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(DBPath string) (*Storage, error) {
	err := initDB(DBPath)
	if err != nil {
		return nil, fmt.Errorf("%e", err)
	}

	db, err := sql.Open("sqlite", DBPath)
	if err != nil {
		return nil, fmt.Errorf("%e", err)
	}

	return &Storage{DB: db}, nil
}

func initDB(DBPath string) error {
	if DBPath == "" {
		appPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("Ошибка вычисления пути: %w", err)
		}
		DBPath = filepath.Join(filepath.Dir(appPath), "scheduler.db")
	}

	if _, err := os.Stat(DBPath); os.IsNotExist(err) {
		if _, err := os.Create(DBPath); err != nil {
			return fmt.Errorf("Ошибка создания файла БД: %w", err)
		}
	}

	db, err := sql.Open("sqlite", DBPath)
	if err != nil {
		return fmt.Errorf("Ошибка подключения к БД: %w", err)
	}
	defer db.Close()

	if _, err = db.Exec("CREATE TABLE IF NOT EXISTS scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, date TEXT, title TEXT, comment TEXT, repeat TEXT)"); err != nil {
		return fmt.Errorf("Ошибка инициализации БД: %w", err)
	}

	return nil
}

func (s *Storage) CreateTask(task models.Task) (int, error) {

	tx, err := s.DB.Begin()
	if err != nil {
		return 0, err
	}

	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat) RETURNING id"
	var id int
	err = tx.QueryRow(
		query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	).Scan(&id)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) GetTasks() ([]models.Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit`

	rows, err := s.DB.Query(
		query,
		sql.Named("limit", models.Limit),
	)
	if err != nil {
		return []models.Task{}, err
	}
	defer rows.Close()

	tasks := make([]models.Task, 0)
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return []models.Task{}, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

func (s *Storage) GetTask(id string) (models.Task, error) {
	var task models.Task

	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	row := s.DB.QueryRow(query, id)
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Task{}, err
		}
		return models.Task{}, err
	}
	return task, nil
}

func (s *Storage) UpdateTask(task models.Task) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	query := "UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id"
	_, err = tx.Exec(
		query,
		sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteTask(id string) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	query := "DELETE FROM scheduler WHERE id = :id"
	_, err = tx.Exec(query, sql.Named("id", id))
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

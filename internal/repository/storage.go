package repository

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(DBPath string) (*sql.DB, error) {
	err := initDB(DBPath)
	if err != nil {
		return nil, fmt.Errorf("%e", err)
	}

	db, err := sql.Open("sqlite", DBPath)
	if err != nil {
		return nil, fmt.Errorf("%e", err)
	}
	
	return db, nil
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

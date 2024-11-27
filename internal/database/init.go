package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func InitDB(dbFile string) error {
	if dbFile == "" {
		appPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("Ошибка вычисления пути: %w", err)
		}
		dbFile = filepath.Join(filepath.Join(filepath.Dir(appPath), "scheduler.db"))
	}

	_, err := os.Stat(dbFile)
	var needSetup bool
	if err != nil {
		needSetup = true
		os.Create(dbFile)
	}

	db, err := sql.Open("sqlite", dbFile)
	defer db.Close()
	if err != nil {
		return fmt.Errorf("Ошибка подключения к БД: %w", err)
	}

	if needSetup {
		query := "CREATE TABLE IF NOT EXISTS scheduler (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `date` CHAR(8) NOT NULL DEFAULT '', `title` VARCHAR(128) NOT NULL DEFAULT '', `comment` VARCHAR(256) NOT NULL DEFAULT '', `repeat` VARCHAR(128) NOT NULL DEFAULT '');"
		_, err = db.Exec(query)
		if err != nil {
			return fmt.Errorf("Ошибка инициализации БД: %w", err)
		}
	}

	return nil
}

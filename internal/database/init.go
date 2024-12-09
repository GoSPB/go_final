package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func InitDB(dbFile string) error {
	// Устанавливаем фиксированный путь к базе данных
	appPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("Ошибка вычисления пути: %w", err)
	}
	dbFile = filepath.Join(filepath.Dir(appPath), "scheduler.db")

	// Проверяем, существует ли файл базы данных
	_, err = os.Stat(dbFile)
	var needSetup bool
	if err != nil {
		needSetup = true
		file, err := os.Create(dbFile)
		if err != nil {
			return fmt.Errorf("Ошибка создания файла БД: %w", err)
		}
		file.Close()
	}

	// Подключаемся к базе данных
	db, err := sql.Open("sqlite", dbFile)
	defer db.Close()
	if err != nil {
		return fmt.Errorf("Ошибка подключения к БД: %w", err)
	}

	// Если требуется настройка, создаём таблицу и индекс
	if needSetup {
		// Создаём таблицу scheduler
		createTableQuery := `
			CREATE TABLE IF NOT EXISTS scheduler (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				date CHAR(8) NOT NULL DEFAULT '',
				title VARCHAR(128) NOT NULL DEFAULT '',
				comment VARCHAR(256) NOT NULL DEFAULT '',
				repeat VARCHAR(128) NOT NULL DEFAULT ''
			);
		`
		_, err = db.Exec(createTableQuery)
		if err != nil {
			return fmt.Errorf("Ошибка создания таблицы: %w", err)
		}

		// Создаём индекс по полю date
		createIndexQuery := "CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);"
		_, err = db.Exec(createIndexQuery)
		if err != nil {
			return fmt.Errorf("Ошибка создания индекса: %w", err)
		}
	}

	return nil
}

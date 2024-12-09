package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/GoSPB/go_final/internal/database"
	"github.com/GoSPB/go_final/internal/models"
	"github.com/GoSPB/go_final/internal/repeat"
)

func UpdateTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		var task models.Task
		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка чтения тела запроса"})
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "Неверный формат запроса"})
			return
		}

		if task.ID == "" {
			json.NewEncoder(w).Encode(map[string]string{"error": "Не указан идентификатор задачи"})
			return
		}

		// Проверяем, что ID — число 32-битное
		if _, err := strconv.ParseInt(task.ID, 10, 32); err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "Неверный идентификатор задачи"})
			return
		}

		if task.Title == "" {
			json.NewEncoder(w).Encode(map[string]string{"error": "Не указан заголовок задачи"})
			return
		}

		now := time.Now()
		if task.Date == "" {
			task.Date = now.Format(models.DateFormat)
		}

		if _, err = time.Parse(models.DateFormat, task.Date); err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "Неверный формат времени"})
			return
		}

		if task.Date < now.Format(models.DateFormat) {
			task.Date = now.Format(models.DateFormat)
		}

		if task.Repeat != "" {
			_, err := repeat.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"error": "Неверный формат правила повтора"})
				return
			}
		}

		err = database.UpdateTask(db, task)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка работы с БД"})
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

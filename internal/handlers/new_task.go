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

func NewTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		now := time.Now()
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

		if task.Title == "" {
			json.NewEncoder(w).Encode(map[string]string{"error": "Не указан заголовок задачи"})
			return
		}

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

		id, err := database.NewTask(db, task)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка работы с БД"})
			return
		}

		idS := strconv.Itoa(int(id))
		json.NewEncoder(w).Encode(map[string]string{"id": idS})
	}
}

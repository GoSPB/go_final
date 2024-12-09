package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/GoSPB/go_final/internal/database"
	"github.com/GoSPB/go_final/internal/repeat"
	"net/http"
	"time"
)

func DoneTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		id := r.URL.Query().Get("id")
		if id == "" {
			json.NewEncoder(w).Encode(map[string]string{"error": "Не указан идентификатор задачи"})
			return
		}

		task, err := database.GetTask(db, id)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Ошибка при получении задачи: %v", err)})
			return
		}

		taskDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Ошибка преобразования даты: %v", err)})
			return
		}

		if task.Repeat != "" {
			nextDateStr, err := repeat.NextDate(taskDate, task.Date, task.Repeat)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Ошибка вычисления следующей даты: %v", err)})
				return
			}

			task.Date = nextDateStr
			err = database.UpdateTask(db, task)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Ошибка при обновлении задачи: %v", err)})
				return
			}
		} else {
			err := database.DeleteTask(db, id)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Ошибка при удалении задачи: %v", err)})
				return
			}
		}

		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/anton-ag/todolist/internal/database"
	rule "github.com/anton-ag/todolist/internal/repeat"
)

func DoneTask(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		id := r.FormValue("id")
		now := time.Now()

		_, err := strconv.Atoi(id)
		if err != nil {
			respondError(w, "Неверный ID")
			return
		}

		task, err := database.GetTask(db, id)
		if err != nil {
			respondError(w, "Задача с данным ID не найдена")
			return
		}

		if task.Repeat == "" {
			err := database.DeleteTask(db, task.ID)
			if err != nil {
				respondError(w, "Ошибка удаления задачи")
				return
			}
			json.NewEncoder(w).Encode(map[string]interface{}{})
			return
		}

		task.Date, err = rule.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			respondError(w, err.Error())
			return
		}
		_, err = database.UpdateTask(db, task)
		if err != nil {
			respondError(w, err.Error())
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

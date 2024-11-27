package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/anton-ag/todolist/internal/database"
)

func GetTask(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		id := r.URL.Query().Get("id")
		if len(id) == 0 {
			respondError(w, "Пустой ID")
			return
		}

		task, err := database.GetTask(db, id)
		if err != nil {
			respondError(w, "Задача с данным ID не найдена")
			return
		}

		body, _ := json.Marshal(task)
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
}

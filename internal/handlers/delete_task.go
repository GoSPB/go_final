package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/anton-ag/todolist/internal/database"
)

func DeleteTask(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		id := r.FormValue("id")

		_, err := strconv.Atoi(id)
		if err != nil {
			respondError(w, "Неверный ID")
			return
		}

		err = database.DeleteTask(db, id)
		if err != nil {
			respondError(w, "Ошибка удаления задачи")
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

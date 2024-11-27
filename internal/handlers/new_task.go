package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/anton-ag/todolist/internal/database"
	"github.com/anton-ag/todolist/internal/models"
	rule "github.com/anton-ag/todolist/internal/repeat"

	_ "modernc.org/sqlite"
)

func NewTask(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		now := time.Now()

		var task models.Task
		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			respondError(w, "Неверный формат запроса")
			return
		}

		if task.Title == "" {
			respondError(w, "Не указан заголовок задачи")
			return
		}

		if task.Date == "" {
			task.Date = now.Format(models.DateFormat)
		}

		if _, err = time.Parse(models.DateFormat, task.Date); err != nil {
			respondError(w, "Неверный формат времени")
			return
		}

		if task.Date < now.Format(models.DateFormat) {
			task.Date = now.Format(models.DateFormat)
		}

		if task.Repeat != "" {
			_, err := rule.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				respondError(w, "Неверный формат правила повтора")
				return
			}
		}

		id, err := database.NewTask(db, task)
		if err != nil {
			respondError(w, "Ошибка работы с БД")
			return
		}

		idS := strconv.Itoa(int(id))
		respondOk(w, idS)
	}
}

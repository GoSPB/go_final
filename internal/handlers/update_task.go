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
)

func UpdateTask(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		now := time.Now()

		var task models.Task
		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			respondError(w, err.Error())
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			respondError(w, "Неверный формат запроса")
			return
		}

		_, err = database.GetTask(db, task.ID)
		if err != nil {
			respondError(w, "Задача с данным ID не найдена")
			return
		}

		_, err = strconv.Atoi(task.ID)
		if err != nil {
			respondError(w, "Неверный ID")
			return
		}

		if task.Title == "" {
			respondError(w, "Не указан заголовок задачи")
			return
		}

		if _, err = time.Parse(models.DateFormat, task.Date); err != nil {
			respondError(w, "Неверный формат времени")
			return
		}

		if task.Repeat != "" {
			_, err := rule.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				respondError(w, "Неверный формат правила повтора")
				return
			}
		}

		id, err := database.UpdateTask(db, task)
		if err != nil {
			respondError(w, err.Error())
			return
		}

		respondOk(w, strconv.Itoa(int(id)))
	}
}

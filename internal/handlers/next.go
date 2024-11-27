package handlers

import (
	"net/http"
	"time"

	"github.com/anton-ag/todolist/internal/models"
	rule "github.com/anton-ag/todolist/internal/repeat"
)

func NextDate(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	nowTime, err := time.Parse(models.DateFormat, now)
	if err != nil {
		respondError(w, "Некорректный формат даты")
		return
	}

	nextDate, err := rule.NextDate(nowTime, date, repeat)
	if err != nil {
		respondError(w, err.Error())
		return
	}

	_, err = w.Write([]byte(nextDate))
	if err != nil {
		respondError(w, err.Error())
		return
	}
}

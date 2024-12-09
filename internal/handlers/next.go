package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/GoSPB/go_final/internal/models"
	rule "github.com/GoSPB/go_final/internal/repeat"
)

func NextDate(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nowTime, err := time.Parse(models.DateFormat, now)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]string{"error": "Неверный формат даты"})
		return
	}

	nextDate, err := rule.NextDate(nowTime, date, repeat)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(nextDate))
}

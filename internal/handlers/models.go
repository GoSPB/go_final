package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/anton-ag/todolist/internal/models"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type IDResponse struct {
	ID string `json:"id"`
}

type TasksResponse struct {
	Tasks []models.Task `json:"tasks"`
}

func respondError(w http.ResponseWriter, s string) {
	body, _ := json.Marshal(ErrorResponse{Error: s})
	w.WriteHeader(http.StatusBadRequest)
	w.Write(body)
}

func respondOk(w http.ResponseWriter, s string) {
	body, _ := json.Marshal(IDResponse{ID: s})
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

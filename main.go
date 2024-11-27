package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/anton-ag/todolist/internal/config"
	"github.com/anton-ag/todolist/internal/database"
	"github.com/anton-ag/todolist/internal/handlers"
	"github.com/go-chi/chi/v5"

	_ "modernc.org/sqlite"
)

func main() {
	var config config.Config
	config.Init()

	err := database.InitDB(config.DBFile)
	if err != nil {
		log.Fatalf("%v", err)
	}
	db, err := sql.Open("sqlite", config.DBFile)
	if err != nil {
		log.Fatalf("Ошибка соединения с БД: %v", err)
	}
	defer db.Close()

	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("web"))
	r.Handle("/*", http.StripPrefix("/", fs))
	r.Get("/api/nextdate", handlers.NextDate)
	r.Post("/api/task", handlers.NewTask(db))
	r.Get("/api/tasks", handlers.GetTasks(db))
	r.Get("/api/task", handlers.GetTask(db))
	r.Put("/api/task", handlers.UpdateTask(db))
	r.Post("/api/task/done", handlers.DoneTask(db))
	r.Delete("/api/task", handlers.DeleteTask(db))

	log.Printf("Запуск сервера на порту %s\n", config.Port)
	err = http.ListenAndServe(config.Port, r)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

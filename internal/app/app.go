package app

import (
	"log"

	"task-traker/internal/config"
	"task-traker/internal/delivery"
	"task-traker/internal/repository"


	_ "modernc.org/sqlite"
)

func Run(cfg *config.Config) {

	storage, err := repository.NewStorage(cfg.DBFile)
	if err != nil {
		log.Fatalf("Ошибка соединения с БД: %v", err)
	}
	defer storage.DB.Close()

	delivery.NewHandler(storage, cfg.Port)
}

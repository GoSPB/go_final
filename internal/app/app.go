package app

import (
	"log"

	"github.com/anton-ag/todolist/internal/config"
	"github.com/anton-ag/todolist/internal/delivery"
	"github.com/anton-ag/todolist/internal/repository"


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

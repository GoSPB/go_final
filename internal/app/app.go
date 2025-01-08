package app



import (
	"log"
	"net/http"

	"github.com/anton-ag/todolist/internal/config"

	"github.com/go-chi/chi/v5"

	_ "modernc.org/sqlite"
)

func Run(cfg *config.Config) {
	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("C:/Users/catas/Downloads/go_final/web"))
	r.Handle("/*", fs)

	log.Printf("Запуск сервера на порту %s\n", cfg.Port)
	err := http.ListenAndServe("localhost:" + cfg.Port, r)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
